package tgClient

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

const (
	tgBotHost         = "api.telegram.org"
	sendMessageMethod = "sendMessage"
)

type TGClient struct {
	host       string
	basePath   string
	httpClient http.Client
}

func New(token string) *TGClient {
	return &TGClient{
		host:       tgBotHost,
		basePath:   newBasePath(token),
		httpClient: http.Client{},
	}
}

func (c *TGClient) Send(username, msg string) error {
	q := url.Values{}
	q.Add("chat_id", username)
	q.Add("text", msg)

	if _, err := c.doRequest(sendMessageMethod, q); err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *TGClient) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme:   "https",
		Host:     c.host,
		Path:     path.Join(c.basePath, method),
		RawQuery: query.Encode(),
	}

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}

	log.Println(string(body))

	return body, nil
}
