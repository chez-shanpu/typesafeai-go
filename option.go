// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of typesafeai-go

package typesafeai

import "net/http"

// Option configures a SystemOneClient.
type Option func(*SystemOneClient)

// WithAPIKey sets the API key sent as a Bearer token in the Authorization
// header. If no key is set, the header is omitted.
func WithAPIKey(key string) Option {
	return func(c *SystemOneClient) {
		c.apiKey = key
	}
}

// WithBaseURL sets the base URL of a TypeSafe AI compatible API, such as
// "http://localhost:8080". The client appends "/v1/systemone" to it.
// Defaults to DefaultBaseURL.
func WithBaseURL(baseURL string) Option {
	return func(c *SystemOneClient) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets the HTTP client used to send requests. A nil client is
// ignored. Defaults to http.DefaultClient.
func WithHTTPClient(client *http.Client) Option {
	return func(c *SystemOneClient) {
		if client != nil {
			c.httpClient = client
		}
	}
}
