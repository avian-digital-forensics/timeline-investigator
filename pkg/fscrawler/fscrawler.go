package fscrawler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

// Client for the fscrawler REST-API
type Client struct {
	httpClient *http.Client
	url        string
	Basepath   string
}

// New creates a new client
// with the specified URL
//
// fs := fscrawler.New("http://localhost:8080/fscrawler")
func New(url string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		url:        url,
	}
}

// Ping to check if fscrawler is healthy
//
// if ok, err := fs.Ping(); !ok {
// 	   if err != nil {
//	       log.Fatal(err)
//	   }
//     log.Fatal("fscrawler not ok")
// }
func (c *Client) Ping(ctx context.Context) (bool, error) {
	resp, err := c.httpClient.Get(c.url)
	if err != nil {
		return false, errors.Wrap(err, "fscrawler.Status: get url")
	}
	defer resp.Body.Close()

	var ping struct {
		OK bool `json:"ok"`
	}

	if err := decodeResponse(resp, &ping); err != nil {
		return false, errors.Wrap(err, "fscrawler.Status: ")
	}

	return ping.OK, nil
}

// Process holds information
// for a process-job
type Process struct {
	client  *Client
	id      string
	index   string
	fileURI string
}

// NewProcess creates a new process to be used for a new processing-job
func (c *Client) NewProcess(fileURI string) *Process {
	return &Process{client: c, fileURI: fileURI}
}

// WithID sets an ID for the processor
func (p *Process) WithID(id string) *Process {
	p.id = id
	return p
}

// WithIndex sets an index name for the processor
func (p *Process) WithIndex(index string) *Process {
	p.index = index
	return p
}

// Start the process-job
func (p *Process) Start(ctx context.Context) error {
	file, err := os.Open(p.fileURI)
	if err != nil {
		return errors.Wrap(err, "cannot open file to process: ")
	}
	defer file.Close()

	// Create a new writer
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	defer writer.Close()

	// Pass custom-ID if it has been added as an argument
	if len(p.id) > 0 {
		if err := writer.WriteField("id", p.id); err != nil {
			return errors.Wrap(err, "cannot set custom-id: ")
		}
	}

	// Pass custom-Index if it has been added as an argument
	if len(p.index) > 0 {
		if err := writer.WriteField("index", p.index); err != nil {
			return errors.Wrap(err, "cannot set custom-index: ")
		}
	}

	// Get part for file
	part, err := writer.CreateFormFile("file", filepath.Base(p.fileURI))
	if err != nil {
		return errors.Wrap(err, "cannot create form-file: ")
	}

	// Write file to part
	if _, err := io.Copy(part, file); err != nil {
		return errors.Wrap(err, "cannot write file to part")
	}

	if err := writer.Close(); err != nil {
		return errors.Wrap(err, "writer.Close(): ")
	}

	// Make the request
	req, err := http.NewRequest("POST", p.client.url+"/_upload", &body)
	if err != nil {
		return errors.Wrap(err, "cannot create new request: ")
	}
	req.WithContext(ctx)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := p.client.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "cannot do request: ")
	}
	defer resp.Body.Close()

	// define the response body
	var respBody struct {
		OK bool `json:"ok"`
	}

	if err := decodeResponse(resp, &respBody); err != nil {
		return err
	}

	if !respBody.OK {
		return fmt.Errorf("Upload was not ok")
	}

	return nil
}

func decodeResponse(r *http.Response, to interface{}) error {
	respBodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "fscrawler.decodeResponse: read response body")
	}

	if err := json.Unmarshal(respBodyBytes, &to); err != nil {
		if r.StatusCode != http.StatusOK {
			return errors.Errorf("fscrawler.decodeResponse: (%d) %s", r.StatusCode, string(respBodyBytes))
		}
		return errors.Wrap(err, "fscrawler.decodeResponse: cannot unmarshal body")
	}
	return nil
}
