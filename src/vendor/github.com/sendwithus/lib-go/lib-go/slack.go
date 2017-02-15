package swu

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// SlackMessage is commented
type SlackMessage struct {
	Text     string `json:"text"`
	Channel  string `json:"channel"`
	Icon     string `json:"icon_emoji"`
	Username string `json:"username"`
}

func slackSendMessage(slackUrl string, msg *SlackMessage) error {
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
