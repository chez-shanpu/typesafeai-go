// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of typesafeai-go

package typesafeai_test

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/chez-shanpu/typesafeai-go"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func testResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestSystemOneDo(t *testing.T) {
	calls := 0
	client := typesafeai.NewSystemOneClient("test-key", &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			calls++
			if req.Method != http.MethodPost {
				t.Errorf("method = %q, want POST", req.Method)
			}
			if req.URL.String() != "https://api.typesafe.ai/v1/systemone" {
				t.Errorf("unexpected URL: %s", req.URL)
			}
			if got := req.Header.Get("Authorization"); got != "Bearer test-key" {
				t.Errorf("Authorization = %q, want Bearer test-key", got)
			}
			if got := req.Header.Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", got)
			}
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read request: %v", err)
			}
			if err := req.Body.Close(); err != nil {
				t.Fatalf("close request: %v", err)
			}
			assertJSONEqual(t, body, `{
				"state":{"text":"sample"},"model":"jev-latest",
				"questions":{
					"relevant":{"type":"noul","instructions":"Relevant?"},
					"category":{"type":"choice","instructions":"Category?","criteria":{"a":null}},
					"quality":{"type":"score","instructions":"Quality?","criteria":["poor","good"]}
				}
			}`)
			return testResponse(http.StatusOK, `{
				"model":"jev-test",
				"answers":{
					"relevant":{"type":"noul","noul":0.8},
					"category":{"type":"choice","choice":"a","probabilities":{"a":1},"confidence":0.9},
					"quality":{"type":"score","score":0.75,"legend":{"0":"poor","1":"good"},
						"probabilities":{"0":0.25,"1":0.75},"confidence":0.8}
				},
				"usage":{"input_tokens":123,"output_tokens":45}
			}`), nil
		}),
	})

	got, err := client.Do(t.Context(), &typesafeai.SystemOneRequest{
		State: map[string]string{"text": "sample"},
		Model: typesafeai.ModelJevLatest,
		Questions: map[string]typesafeai.Question{
			"relevant": &typesafeai.NoulQuestion{Instructions: "Relevant?"},
			"category": &typesafeai.ChoiceQuestion{
				Instructions: "Category?", Criteria: map[string]*string{"a": nil},
			},
			"quality": &typesafeai.ScoreQuestion{
				Instructions: "Quality?", Criteria: []string{"poor", "good"},
			},
		},
	})
	if err != nil {
		t.Fatalf("Do() error: %v", err)
	}
	if calls != 1 {
		t.Errorf("HTTP calls = %d, want 1", calls)
	}
	if got.Model != "jev-test" {
		t.Errorf("Model = %q, want jev-test", got.Model)
	}
	if want := (typesafeai.TokenUsage{InputTokens: 123, OutputTokens: 45}); got.Usage != want {
		t.Errorf("Usage = %+v, want %+v", got.Usage, want)
	}
	wantAnswers := typesafeai.QuestionKeyToAnswer{
		"relevant": &typesafeai.NoulAnswer{Noul: 0.8},
		"category": &typesafeai.ChoiceAnswer{
			Choice: "a", Probabilities: map[string]float64{"a": 1}, Confidence: 0.9,
		},
		"quality": &typesafeai.ScoreAnswer{
			Score: 0.75, Legend: map[string]string{"0": "poor", "1": "good"},
			Probabilities: map[string]float64{"0": 0.25, "1": 0.75}, Confidence: 0.8,
		},
	}
	if !reflect.DeepEqual(got.Answers, wantAnswers) {
		t.Errorf("Answers = %#v, want %#v", got.Answers, wantAnswers)
	}
}

func TestSystemOneDoAPIError(t *testing.T) {
	for _, status := range []int{199, 300, 400, 401, 429, 500} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			client := typesafeai.NewSystemOneClient("test-key", &http.Client{
				Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
					resp := testResponse(status, "not JSON")
					resp.Header.Set("x-typesafe-request-id", "req-123")
					return resp, nil
				}),
			})

			got, err := client.Do(t.Context(), &typesafeai.SystemOneRequest{})

			apiErr, ok := errors.AsType[*typesafeai.APIError](err)
			if !ok {
				t.Fatalf("error = %v, want APIError", err)
			}
			if apiErr.StatusCode != status {
				t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, status)
			}
			if apiErr.RequestID != "req-123" {
				t.Errorf("RequestID = %q, want req-123", apiErr.RequestID)
			}
			if string(apiErr.Body) != "not JSON" {
				t.Errorf("Body = %q, want %q", apiErr.Body, "not JSON")
			}
			if got != nil {
				t.Errorf("response = %#v, want nil", got)
			}
		})
	}
}

func TestSystemOneDoConnectionError(t *testing.T) {
	wantErr := errors.New("connection failed")
	client := typesafeai.NewSystemOneClient("test-key", &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, wantErr
		}),
	})
	got, err := client.Do(t.Context(), &typesafeai.SystemOneRequest{})
	if _, ok := errors.AsType[*typesafeai.ConnectionError](err); !ok {
		t.Fatalf("error = %v, want ConnectionError", err)
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want wrapped %v", err, wantErr)
	}
	if got != nil {
		t.Errorf("response = %#v, want nil", got)
	}
}
