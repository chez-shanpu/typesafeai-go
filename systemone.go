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
)

// SystemOneV1URL is the System One v1 API endpoint.
const SystemOneV1URL = "https://api.typesafe.ai/v1/systemone"

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
	httpClient *http.Client
}

// NewSystemOneClient returns a SystemOneClient that authenticates with
// apiKey and issues requests via client.
func NewSystemOneClient(apiKey string, client *http.Client) *SystemOneClient {
	return &SystemOneClient{
		apiKey:     apiKey,
		httpClient: client,
	}
}

// Do sends req to the System One API and returns the decoded response.
// It returns a *ConnectionError on network failure, an *APIError for
// non-2xx responses, or a *ResponseValidationError if the response body
// cannot be decoded.
func (c *SystemOneClient) Do(ctx context.Context, req *SystemOneRequest) (sresp *SystemOneResponse, err error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, SystemOneV1URL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	r.Header.Set("Authorization", "Bearer "+c.apiKey)
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
