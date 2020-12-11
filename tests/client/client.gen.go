// Code generated by oto; DO NOT EDIT.

package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	context "context"
	"github.com/pkg/errors"
	http "net/http"
)

// Client is used to access Pace services.
type Client struct {
	// RemoteHost is the URL of the remote server that this Client should
	// access.
	RemoteHost string
	// HTTPClient is the http.Client to use when making HTTP requests.
	HTTPClient *http.Client
	// Debug writes a line of debug log output.
	Debug func(s string)
}

// New makes a new Client.
func New(remoteHost string) *Client {
	c := &Client{
		RemoteHost: remoteHost,
		Debug:      func(s string) {},
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
	return c
}

// Service is the main-service
type Service struct {
	client *Client
	token  string
}

// NewService makes a new client for accessing Service services.
func NewService(client *Client, token string) *Service {
	return &Service{
		client: client,
		token:  token,
	}
}

// Greet sends a polite greeting
func (s *Service) Greet(ctx context.Context, r GreetRequest) (*GreetResponse, error) {
	requestBodyBytes, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "Service.Greet: marshal GreetRequest")
	}
	url := s.client.RemoteHost + "Service.Greet"
	s.client.Debug(fmt.Sprintf("POST %s", url))
	s.client.Debug(fmt.Sprintf(">> %s", string(requestBodyBytes)))
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestBodyBytes))
	if err != nil {
		return nil, errors.Wrap(err, "Service.Greet: NewRequest")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Authorization", s.token)
	req = req.WithContext(ctx)
	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Service.Greet")
	}
	defer resp.Body.Close()
	var response struct {
		GreetResponse
		Error string
	}
	var bodyReader io.Reader = resp.Body
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		decodedBody, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Service.Greet: new gzip reader")
		}
		defer decodedBody.Close()
		bodyReader = decodedBody
	}
	respBodyBytes, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "Service.Greet: read response body")
	}
	s.client.Debug(fmt.Sprintf("<< %s", string(respBodyBytes)))
	if err := json.Unmarshal(respBodyBytes, &response); err != nil {
		if resp.StatusCode != http.StatusOK {
			return nil, errors.Errorf("Service.Greet: (%d) %v", resp.StatusCode, string(respBodyBytes))
		}
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return &response.GreetResponse, nil
}

// GreetRequest is the request object for GreeterService.Greet.
type GreetRequest struct {
	// Namee of the person to greet
	Name string `json:"name"`
}

// GreetResponse is the response object containing a person's greeting.
type GreetResponse struct {
	// Greeting is a nice message welcoming somebody.
	Greeting string `json:"greeting"`
}
