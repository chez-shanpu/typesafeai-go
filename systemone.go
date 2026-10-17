// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of typesafeai-go

package typesafeai

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// DefaultBaseURL is the base URL of the official TypeSafe AI API.
const DefaultBaseURL = "https://api.typesafe.ai"

// Question type discriminators used in the "type" JSON field.
const (
	TypeNoul   = "noul"
	TypeChoice = "choice"
	TypeScore  = "score"
)

// SystemOneRequest is the request body sent to the System One API.
type SystemOneRequest struct {
	State     any                 `json:"state"`
	Model     string              `json:"model"`
	Questions map[string]Question `json:"questions"`
}

// SystemOneResponse is the response body returned by the System One API.
type SystemOneResponse struct {
	Model   string              `json:"model"`
	Answers QuestionKeyToAnswer `json:"answers"`
	Usage   TokenUsage          `json:"usage"`
}

// TokenUsage reports token consumption for a request.
type TokenUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// SystemOneClient calls the System One API.
type SystemOneClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewSystemOneClient returns a SystemOneClient configured by opts.
// Without options, it sends unauthenticated requests to DefaultBaseURL
// via http.DefaultClient.
func NewSystemOneClient(opts ...Option) *SystemOneClient {
	c := &SystemOneClient{
		baseURL:    DefaultBaseURL,
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// endpoint returns the System One v1 endpoint under the configured base URL.
func (c *SystemOneClient) endpoint() (string, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", &ValidationError{Field: "baseURL", Message: err.Error()}
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return "", &ValidationError{
			Field:   "baseURL",
			Message: fmt.Sprintf("must be an absolute http(s) URL, got %q", c.baseURL),
		}
	}
	return u.JoinPath("v1", "systemone").String(), nil
}

// Do sends req to the System One API and returns the decoded response.
// It returns a *ValidationError if the base URL is invalid, a
// *ConnectionError on network failure, an *APIError for
// non-2xx responses, or the error from encoding/json if the response
// body cannot be decoded.
func (c *SystemOneClient) Do(ctx context.Context, req *SystemOneRequest) (sresp *SystemOneResponse, err error) {
	endpoint, err := c.endpoint()
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	if c.apiKey != "" {
		r.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	r.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, &ConnectionError{Err: err}
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close response body: %w", cerr)
		}
	}()

	reqID := resp.Header.Get("x-typesafe-request-id")
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &ConnectionError{Err: err}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			RequestID:  reqID,
			Body:       respBody,
		}
	}

	// TODO validate response body against schema
	if err = json.Unmarshal(respBody, &sresp); err != nil {
		return nil, err
	}
	return sresp, nil
}
