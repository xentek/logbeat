package logbeat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// ISO8601 Date Format.
const ISO8601 = "2006-01-02T15:04:05Z07:00"

// OpbeatClient is used to communicate with Opbeat's API.
type OpbeatClient struct {
	Endpoint string
	Http     *http.Client
	Token    string
}

// OpbeatPayload structures log entries for Opbeat's API.
type OpbeatPayload struct {
	Extra     OpbeatExtra   `json:"extra,omitempty"`
	Level     string        `json:"level"`
	Logger    string        `json:"logger"`
	Machine   OpbeatMachine `json:"machine,omitempty"`
	Message   string        `json:"message"`
	Timestamp string        `json:"timestamp"`
}

// OpbeatExtra structures Logrus Fields for Opbeat's API.
type OpbeatExtra map[string]interface{}

// OpbeatMachine represents the hostname of the Machine the error occured on to Opbeat's API.
type OpbeatMachine struct {
	Hostname string `json:"hostname"`
}

// NewOpbeatClient returns an OpbeatClient used for commnicating with Opbeat's API.
func NewOpbeatClient(org, app, token string) *OpbeatClient {
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	return &OpbeatClient{
		Endpoint: OpbeatEndpoint(org, app),
		Http:     client,
		Token:    token,
	}
}

// NewOpbeatMachine returns an OpbeatMachine for the current machine.
func NewOpbeatMachine() OpbeatMachine {
	var machine = OpbeatMachine{}
	hostname, err := os.Hostname()
	if err != nil {
		return machine
	}
	machine.Hostname = hostname
	return machine
}

// OpbeatLevel returns the logrus.Level, as a string that Opbeat will accept, for the given logrus.Entry.
func OpbeatLevel(entry *logrus.Entry) string {
	level := entry.Level.String()
	if level == "panic" {
		level = "critical"
	}
	return level
}

// NewOpbeatExtra returns an OpbeatExtra for the given logrus.Entry.
func NewOpbeatExtra(entry *logrus.Entry) OpbeatExtra {
	var extra = OpbeatExtra{}
	for k, v := range entry.Data {
		extra[k] = v
	}
	return extra
}

// NewOpbeatPayload returns an OpbeatPayload for the given logrus.Entry.
func NewOpbeatPayload(entry *logrus.Entry) *OpbeatPayload {
	return &OpbeatPayload{
		Extra:     NewOpbeatExtra(entry),
		Level:     OpbeatLevel(entry),
		Logger:    fmt.Sprintf("logbeat-%s", LogbeatVersion),
		Machine:   NewOpbeatMachine(),
		Message:   entry.Message,
		Timestamp: entry.Time.UTC().Format(ISO8601),
	}
}

// JSON returns the OpbeatPayload, seralized as JSON. It returns a
// bytes.Buffer to satisfy the io.Reader interface required by
// http.Request.
func (payload *OpbeatPayload) JSON() (*bytes.Buffer, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonPayload), nil
}

// NewOpbeatRequest creates an http.Request used to notify the Opbeat API.
func (client *OpbeatClient) NewOpbeatRequest(json *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest("POST", client.Endpoint, json)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", OpbeatBearerAuth(client.Token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", OpbeatUserAgent())

	return req, nil
}

// Notify sends a JSON encoded OpbeatPayload to Opbeat's API.
func (client *OpbeatClient) Notify(entry *logrus.Entry) (*http.Response, error) {
	payload := NewOpbeatPayload(entry)

	body, err := payload.JSON()
	if err != nil {
		return nil, err
	}

	req, err := client.NewOpbeatRequest(body)
	if err != nil {
		return nil, err
	}

	return client.Http.Do(req)
}

// OpbeatUserAgent to identify the LogbeatHook to Opbeat's API.
func OpbeatUserAgent() string {
	return fmt.Sprintf("Logbeat/%s (+https://github.com/xentek/logbeat)", LogbeatVersion)
}

// OpbeatEndpoint returns a formatted URI to Opbeat's API.
func OpbeatEndpoint(org, app string) string {
	return fmt.Sprintf("https://intake.opbeat.com/api/v1/organizations/%s/apps/%s/errors/", org, app)
}

// OpbeatBearerAuth Formats the Authorization header value for Opbeat's API.
func OpbeatBearerAuth(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}
