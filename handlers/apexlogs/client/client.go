package client

import (
	"fmt"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// TODO: generate full client
// TODO: docs
// TODO: auth

// Error is an error returned by the client.
type Error struct {
	Method     string
	Status     string
	StatusCode int
	Type       string
	Message    string
}

// Error implementation.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Event represents a single log event.
type Event struct {
	// ID is the event id.
	ID string `json:"id"`

	// Level is the severity level.
	Level string `json:"level"`

	// Message is the log message.
	Message string `json:"message"`

	// Fields is the log fields.
	Fields map[string]interface{} `json:"fields"`

	// Timestamp is the creation timestamp.
	Timestamp time.Time `json:"timestamp"`
}

// AddEventsInput params.
type AddEventsInput struct {
	// ProjectID is the project id.
	ProjectID string `json:"project_id"`

	// Events is the batch of events.
	Events []Event `json:"events"`
}

// Client is the API client.
type Client struct {
	URL string
	AuthToken string
}

// AddEvents ingested a batch of events.
func (c *Client) AddEvents(in AddEventsInput) error {
	return c.call("", "add_events", nil, in)
}

// call invokes a method on the given service.
func (c *Client) call(service, name string, out interface{}, in ...interface{}) error {
	var body io.Reader

	// input params
	if len(in) > 0 {
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(in)
		if err != nil {
			return err
		}
		body = &buf
	}

	// endpoint
	url := c.URL + "/" + service + "/" + name
	if service == "" {
		url = c.URL + "/" + name
	}

	// POST request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// auth token
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	// response
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		var e Error
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		}
		e.Status = res.Status
		e.StatusCode = res.StatusCode
		e.Method = name
		return e
	}

	// output params
	if out != nil {
		err = json.NewDecoder(res.Body).Decode(out)
		if err != nil {
			return err
		}
	}

	return nil
}
