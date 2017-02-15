package swu

import (
	"encoding/base64"
	"errors"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"strings"
)

var swuMsgKey = "SWU-MESSAGE-KEY"

// Using space because it is JSON encoded
const paddingValue = byte(' ')

type SQSService interface {
	SendMessage(queueName *string, message *[]byte) (*string, error)
	ReceiveMessages(queueName *string, maxMessages int64, waitTimeSeconds int64) ([]SQSMessage, error)
	DeleteMessage(queueName *string, message SQSMessage) error
	GetQueueSize(queueName string, includeInFlight bool) (int64, error)
	QueuesEmpty(queueNames []string, includeInFlight bool) (bool, error)
	GetQueueSizes(prefix string) ([]QueueInfo, error)
}

type sqsServiceImpl struct {
	SQS *sqs.SQS
}

var cachedService = sqsServiceImpl{
	SQS: sqs.New(getAwsSession()),
}

func NewSQSService() SQSService {
	return sqsServiceImpl{
		SQS: sqs.New(getAwsSession()),
	}
}

func NewSharedSQSService() SQSService {
	return &cachedService
}

type SQSMessage struct {
	Id            string
	Body          string
	ReceiptHandle string
}

type QueueInfo struct {
	Name string
	Size int64
}

func (s sqsServiceImpl) SendMessage(queueName *string, message *[]byte) (*string, error) {

	// Check if we have a valid queue
	if queueName == nil {
		return nil, errors.New("Cannot write to queue ``")
	}

	// Encrypt the data before placing it into the message
	encrypted_message, public_key, err := EncryptAES(message, paddingValue)
	if err != nil {
		panic(err)
	}

	// convert encrypted data and public key to base64 encoding
	b64_encrypted_message := base64.StdEncoding.EncodeToString(*encrypted_message)
	b64_public_key := base64.StdEncoding.EncodeToString(*public_key)

	// store the public key in the sqs message so the data can be decrypted on the other end.
	attrs := map[string]*sqs.MessageAttributeValue{
		swuMsgKey: {
			DataType:    aws.String("String"),
			StringValue: aws.String(string(b64_public_key)),
		},
	}

	// Grab the SQS queue and place the message on it. Return error if any.
	url, err := s.getQueueURL(*queueName)
	if err != nil {
		return nil, err
	}
	params := &sqs.SendMessageInput{
		QueueUrl:          aws.String(*url),
		MessageBody:       aws.String(string(b64_encrypted_message)),
		MessageAttributes: attrs,
	}

	output, err := s.SQS.SendMessage(params) // Could do something with the response

	return output.MessageId, err
}

func (s sqsServiceImpl) ReceiveMessages(queueName *string, maxMessages int64, waitTimeSeconds int64) ([]SQSMessage, error) {
	attrs := []*string{&swuMsgKey}

	url, err := s.getQueueURL(*queueName)
	if err != nil {
		return nil, err
	}
	params := &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(*url),
		MaxNumberOfMessages:   aws.Int64(maxMessages),
		WaitTimeSeconds:       aws.Int64(waitTimeSeconds),
		MessageAttributeNames: attrs,
	}

	resp, err := s.SQS.ReceiveMessage(params)
	if err != nil {
		return nil, err
	}

	messages := make([]SQSMessage, len(resp.Messages))
	for index, message := range resp.Messages {
		b64_public_key := message.MessageAttributes[swuMsgKey].StringValue
		b64_body := message.Body

		public_key, err := base64.StdEncoding.DecodeString(*b64_public_key)
		if err != nil {
			panic(err)
		}
		encrypted_body, err := base64.StdEncoding.DecodeString(*b64_body)
		if err != nil {
			panic(err)
		}
		decrypted_body, err := DecryptAES(&encrypted_body, &public_key, paddingValue)
		if err != nil {
			panic(err)
		}

		messages[index] = SQSMessage{
			Id:            *message.MessageId,
			Body:          string(*decrypted_body),
			ReceiptHandle: *message.ReceiptHandle,
		}

	}

	return messages, nil
}

func (s sqsServiceImpl) DeleteMessage(queueName *string, message SQSMessage) error {

	url, err := s.getQueueURL(*queueName)
	if err != nil {
		return err
	}

	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(*url),
		ReceiptHandle: aws.String(message.ReceiptHandle),
	}

	_, err = s.SQS.DeleteMessage(params)
	if err != nil {
		return err
	}

	return nil
}

func (s sqsServiceImpl) getQueueIntAttributes(queueUrl string, attributeNames ...string) (map[string]int64, error) {
	searchNames := make([]*string, len(attributeNames))
	for index, name := range attributeNames {
		searchNames[index] = aws.String(name)
	}

	params := &sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(queueUrl), // Required
		AttributeNames: searchNames,
	}

	response, err := s.SQS.GetQueueAttributes(params)
	if err != nil {
		return nil, err
	}

	values := make(map[string]int64, len(attributeNames))
	for attrName, attrValue := range response.Attributes {
		attrValueInt, err := strconv.ParseInt(*attrValue, 10, 64)
		values[attrName] = attrValueInt
		if err != nil {
			return values, err
		}
	}
	return values, nil
}

func (s sqsServiceImpl) GetQueueSize(queueName string, includeInFlight bool) (int64, error) {
	url, err := s.getQueueURL(queueName)
	if err != nil {
		return -1, err
	}
	return s.GetQueueSizeByUrl(url, true)
}

func (s sqsServiceImpl) GetQueueSizeByUrl(url *string, includeInFlight bool) (int64, error) {
	values, err := s.getQueueIntAttributes(*url, "ApproximateNumberOfMessages", "ApproximateNumberOfMessagesNotVisible")
	if err != nil {
		return -1, err // using -1 here because there was an error, and this makes it obvious that we don't know the queue size
	}
	size := values["ApproximateNumberOfMessages"]
	if includeInFlight {
		inFlightSize := values["ApproximateNumberOfMessagesNotVisible"]
		size += inFlightSize
	}
	return size, nil
}

func (s sqsServiceImpl) QueuesEmpty(queueNames []string, includeInFlight bool) (bool, error) {
	for _, queueName := range queueNames {
		size, err := s.GetQueueSize(queueName, includeInFlight)
		empty := size == 0
		if !empty || err != nil {
			return false, err
		}
	}
	return true, nil
}

func (s sqsServiceImpl) GetQueueSizes(prefix string) ([]QueueInfo, error) {

	input := sqs.ListQueuesInput{
		QueueNamePrefix: aws.String(prefix),
	}

	result, err := s.SQS.ListQueues(&input)
	if err != nil {
		return nil, err
	}

	items := []QueueInfo{}
	for _, url := range result.QueueUrls {
		size, err := s.GetQueueSizeByUrl(url, true)
		if err != nil {
			return items, err
		}

		// the queue name is always the last bit of the queue url
		items = append(items, QueueInfo{
			Name: (*url)[strings.LastIndex(*url, "/")+1:],
			Size: size,
		})
	}
	return items, nil
}

// TODO: We should really cache QueueUrls so we're not looking them up on every request
func (s sqsServiceImpl) getQueueURL(queueName string) (*string, error) {

	params := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}

	response, err := s.SQS.GetQueueUrl(params)
	if err != nil {
		return nil, err
	}

	return response.QueueUrl, nil
}
