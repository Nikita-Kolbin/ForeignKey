package notification_client

import (
	"ForeignKey/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type NotificationClient struct {
	client      http.Client
	emailURL    string
	telegramURL string
}

func New(cfg config.NotificationClientConfig) *NotificationClient {
	return &NotificationClient{
		client:      http.Client{},
		emailURL:    cfg.EmailURL,
		telegramURL: cfg.TelegramURL,
	}
}

func (nc *NotificationClient) SendEmail(email, msg string) error {
	const op = "notification_client.SendEmail"

	n := EmailNotification{
		ReceiverEmail: email,
		Message:       msg,
	}

	resp, err := nc.doRequest(nc.emailURL, n)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = nc.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (nc *NotificationClient) SendTelegram(username, msg string) error {
	const op = "notification_client.SendTelegram"

	n := TelegramNotification{
		Username: username,
		Message:  msg,
	}

	resp, err := nc.doRequest(nc.telegramURL, n)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = nc.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (nc *NotificationClient) handleResponse(resp *Response) error {
	switch resp.Status {
	case statusOk:
		return nil
	case statusError:
		return fmt.Errorf("%s", resp.Error)
	default:
		return ErrUnknownStatus
	}
}

func (nc *NotificationClient) doRequest(url string, body interface{}) (*Response, error) {
	j, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(j)

	r, err := nc.client.Post(
		url,
		contentJSON,
		reader,
	)
	if err != nil {
		return nil, err
	}

	byteResp, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var resp Response
	err = json.Unmarshal(byteResp, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, err
}
