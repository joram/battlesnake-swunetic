package swu

import (
	"testing"
)

type TestSlackServiceImp struct {
	actual  slackServiceImp
	lastMsg *SlackMessage
}
func (s *TestSlackServiceImp) channelExists(channelName string, existingChannelNames []string) bool {
	return s.actual.channelExists(channelName, existingChannelNames)
}
func (s *TestSlackServiceImp) getSlackChannels(token string) ([]string, error) {
	return []string{"valid_channel_A", "valid_channel_B", "fires"}, nil
}
func (s *TestSlackServiceImp) sendMessage(slackUrl string, msg *SlackMessage) error {
	s.lastMsg = msg
	return nil
}

func TestSlackChannelExists(t *testing.T) {
	slackService := TestSlackServiceImp{
		actual:slackServiceImp{},
	}
	existingChannels := []string{"valid_channel_A", "valid_channel_B", "fires"}

	if slackService.channelExists("#a_test_channel_that_doesnt_exist", existingChannels) {
		t.Error("Channel non-existance check failed")
	}

	if !slackService.channelExists("#fires", existingChannels) {
		t.Error("Channel existance check failed")
	}

	logger := NewLogger("")
	config := SlackConfig{"url", "username", "icon", "token"}
	logger.slackService = &slackService
	logger.InitializeSlack(&config)

	// Act
	logger.Slack("#valid_channel_B", "lude message 1")

	if slackService.lastMsg == nil {
		t.Error("didn't send the message, we should have")
	}
	if slackService.lastMsg.Channel != "#valid_channel_B" {
		t.Error("Channel changed to when it shouldn't have: ", slackService.lastMsg.Channel)
	}
	if slackService.lastMsg.Text != "lude message 1" {
		t.Error("Text changed to when it shouldn't have: ", slackService.lastMsg.Text)
	}

	slackService.lastMsg = nil

	// Act
	logger.Slack("#invalid_channel_C", "lude message 2")

	if slackService.lastMsg == nil {
		t.Error("didn't send the message, we should have")
	}
	if slackService.lastMsg.Channel != "#fires" {
		t.Error("Channel changed to when it shouldn't have: ", slackService.lastMsg.Channel)
	}
	if slackService.lastMsg.Text != "Slack channel does not exist `#invalid_channel_C`: lude message 2" {
		t.Error("Text changed to when it shouldn't have: ", slackService.lastMsg.Text)
	}


	// TODO test without a token
}