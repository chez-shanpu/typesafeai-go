// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of typesafeai-go

package typesafeai

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
)

// Question is implemented by NoulQuestion, ChoiceQuestion, and ScoreQuestion.
type Question interface {
	isQuestion()
}

// NoulQuestion asks the model a true/false-style question.
type NoulQuestion struct {
	Instructions any           `json:"instructions"`
	Criteria     *NoulCriteria `json:"criteria,omitempty"`
}

// NoulCriteria describes what the true and false outcomes of a NoulQuestion mean.
type NoulCriteria struct {
	True  *string `json:"true,omitempty"`
	False *string `json:"false,omitempty"`
}

func (n *NoulQuestion) MarshalJSONTo(enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, struct {
		Type         string        `json:"type"`
		Instructions any           `json:"instructions"`
		Criteria     *NoulCriteria `json:"criteria,omitempty"`
	}{
		Type:         TypeNoul,
		Instructions: n.Instructions,
		Criteria:     n.Criteria,
	})
}

func (*NoulQuestion) isQuestion() {}

// ChoiceQuestion asks the model to select one of several named choices.
type ChoiceQuestion struct {
	Instructions any                `json:"instructions"`
	Criteria     map[string]*string `json:"criteria"`
}

func (c *ChoiceQuestion) MarshalJSONTo(enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, struct {
		Type         string             `json:"type"`
		Instructions any                `json:"instructions"`
		Criteria     map[string]*string `json:"criteria"`
	}{
		Type:         TypeChoice,
		Instructions: c.Instructions,
		Criteria:     c.Criteria,
	})
}

func (*ChoiceQuestion) isQuestion() {}

// ScoreQuestion asks the model to produce a numeric score.
type ScoreQuestion struct {
	Instructions any      `json:"instructions"`
	Criteria     []string `json:"criteria"`
}

func (s *ScoreQuestion) MarshalJSONTo(enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, struct {
		Type         string   `json:"type"`
		Instructions any      `json:"instructions"`
		Criteria     []string `json:"criteria"`
	}{
		Type:         TypeScore,
		Instructions: s.Instructions,
		Criteria:     s.Criteria,
	})
}

func (*ScoreQuestion) isQuestion() {}
