package rocketchat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	url        string
	userID     string
	authToken  string
}

// Message defines a chat message.
type Message struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
	Alias   string `json:"alias,omitempty"`
	Emoji   string `json:"emoji,omitempty"`
}

// User represents a user in Rocket.Chat.
type User struct {
	ID       string `json:"_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// Channel represents a channel in Rocket.Chat.
type Channel struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

// NewClient returns a new Rocket.Chat client.
func NewClient(url, userID, authToken string) *Client {
	return &Client{
		&http.Client{
			Timeout: 30 * time.Second,
		},
		url,
		userID,
		authToken,
	}
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", c.url, path)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-User-Id", c.userID)
	req.Header.Set("X-Auth-Token", c.authToken)

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		errorResponse := struct {
			Error string `json:"error"`
		}{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return fmt.Errorf("request to '%s' failed with http code %d", resp.Request.URL, resp.StatusCode)
		}
		return fmt.Errorf("request to '%s' failed with http code %d: %s", resp.Request.URL, resp.StatusCode, errorResponse.Error)
	}

	// nil for empty responses
	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return err
}

func (c *Client) SendMessage(msg *Message) error {
	req, err := c.newRequest("POST", "/api/v1/chat.postMessage", msg)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

func (c *Client) ListChannels() ([]Channel, error) {
	channelsResponse := struct {
		Channels []Channel `json:"channels"`
		Count    int       `json:"count"`
	}{}

	req, err := c.newRequest("GET", "/api/v1/channels.list", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("count", "500")
	req.URL.RawQuery = q.Encode()

	err = c.do(req, &channelsResponse)
	if err != nil {
		return nil, err
	}

	return channelsResponse.Channels, nil
}

func (c *Client) ListUsers() ([]User, error) {
	userResponse := struct {
		Users []User `json:"users"`
		Count int    `json:"count"`
	}{}

	req, err := c.newRequest("GET", "/api/v1/users.list", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("count", "500")
	req.URL.RawQuery = q.Encode()

	err = c.do(req, &userResponse)
	if err != nil {
		return nil, err
	}

	return userResponse.Users, nil
}

func (c *Client) TestConnection() error {
	req, err := c.newRequest("GET", "/api/v1/me", nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
