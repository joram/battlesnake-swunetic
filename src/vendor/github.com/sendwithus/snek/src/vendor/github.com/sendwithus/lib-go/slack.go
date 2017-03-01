package swu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
)

var slackToken string

type SlackConfig struct {
	Url      string
	Username string
	Icon     string
	Token    string
}

// SlackMessage is commented
type SlackMessage struct {
	Text     string `json:"text"`
	Channel  string `json:"channel"`
	Icon     string `json:"icon_emoji"`
	Username string `json:"username"`
}

type SlackChannelDetails struct {
	Name        string `json:"name"`
	Is_channel  bool   `json:"is_channel"`
	Is_archived bool   `json:"is_archived"`
}

type ChannelListResponse struct {
	Ok       bool
	Channels []SlackChannelDetails
}

type slackService interface {
	sendMessage(slackUrl string, msg *SlackMessage) error
	channelExists(channelName string, existingChannelNames []string) bool
	getSlackChannels(token string) ([]string, error)
}

type slackServiceImp struct {}
func NewSlackService() slackService {
	return &slackServiceImp{}
}

func (s *slackServiceImp) sendMessage(slackUrl string, msg *SlackMessage) error {
	if slackUrl == "" {
		return nil
	}

	// Post to Slack
	data, err := json.Marshal(*msg)
	if err != nil {
		return err
	}

	response, err := http.Post(slackUrl, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return nil
}

func (s *slackServiceImp) channelExists(channelName string, existingChannelNames []string) bool {
	channelName = strings.Replace(channelName, "#", "", 1)
	for _, existingChannelName := range existingChannelNames {
		if existingChannelName == channelName {
			return true
		}
	}

	return false
}


func (s *slackServiceImp) getSlackChannels(token string) ([]string, error) {
	response, err := http.Get(fmt.Sprintf("https://slack.com/api/channels.list?token=%v&pretty=1", token))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var r ChannelListResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	var channels []string
	for _, channelDetails := range r.Channels {
		if !channelDetails.Is_archived {
			channels = append(channels, channelDetails.Name)
		}
	}

	return channels, nil
}
