package rocketchat

import (
	"bytes"
	"context"
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
	Text string `json:"text"`

	// Channel can either be a channel prefixed with a # or a username prefixed with an @.
	Channel string `json:"channel"`

	// Name under which the message appears.
	Alias string `json:"alias,omitempty"`
	Emoji string `json:"emoji,omitempty"`
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

func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", c.url, path)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, buf)
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

// do performs a http request. If result is not nil it tries to json decode the
// body into the result.
func (c *Client) do(req *http.Request, result any) error {
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

	if result != nil {
		err = json.NewDecoder(resp.Body).Decode(result)
	}
	return err
}

// SendMessage sends a message to a channel or user.
func (c *Client) SendMessage(ctx context.Context, msg *Message) error {
	req, err := c.newRequest(ctx, "POST", "/api/v1/chat.postMessage", msg)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

// ListChannels returns a list of channels.
func (c *Client) ListChannels(ctx context.Context) ([]Channel, error) {
	channelsResponse := struct {
		Channels []Channel `json:"channels"`
		Count    int       `json:"count"`
	}{}

	req, err := c.newRequest(ctx, "GET", "/api/v1/channels.list", nil)
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

// ListUsers returns a list of users.
func (c *Client) ListUsers(ctx context.Context) ([]User, error) {
	userResponse := struct {
		Users []User `json:"users"`
		Count int    `json:"count"`
	}{}

	req, err := c.newRequest(ctx, "GET", "/api/v1/users.list", nil)
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

// TestConnection calls the endpoint /api/v1/me. It is intended to verify if the
// connection to the Rocket.Chat server works.
func (c *Client) TestConnection(ctx context.Context) error {
	req, err := c.newRequest(ctx, "GET", "/api/v1/me", nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
