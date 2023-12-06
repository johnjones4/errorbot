package errorbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type outgoingMessage struct {
	Text   string `json:"text"`
	ChatId int    `json:"chat_id"`
}

type telegram struct {
	token string
}

func (t *telegram) callMethod(name string, parameters interface{}, response interface{}) error {
	bodyBytes := []byte{}
	var err error

	if parameters != nil {
		bodyBytes, err = json.Marshal(parameters)
		if err != nil {
			return err
		}
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", t.token, name)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response: %d %s", resp.StatusCode, string(responseBody))
	}

	if response == nil {
		return nil
	}

	err = json.Unmarshal(responseBody, response)
	if err != nil {
		return err
	}

	return nil
}

func (t *telegram) sendMessage(m outgoingMessage) error {
	err := t.callMethod("sendMessage", m, nil)
	if err != nil {
		return err
	}
	return nil
}
