package typesafeai

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
)

// Answer is implemented by NoulAnswer, ChoiceAnswer, and ScoreAnswer.
type Answer interface {
	isAnswer()
}

// QuestionKeyToAnswer maps a request question key to its answer, decoding
// each value into the concrete Answer type indicated by its "type" field.
type QuestionKeyToAnswer map[string]Answer

func (q *QuestionKeyToAnswer) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	var raw map[string]jsontext.Value
	if err := json.UnmarshalDecode(dec, &raw); err != nil {
		return err
	}

	out := make(QuestionKeyToAnswer, len(raw))
	for key, v := range raw {
		var disc struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(v, &disc); err != nil {
			return err
		}

		var a Answer
		switch disc.Type {
		case TypeNoul:
			a = new(NoulAnswer)
		case TypeChoice:
			a = new(ChoiceAnswer)
		case TypeScore:
			a = new(ScoreAnswer)
		default:
			return fmt.Errorf("unknown answer type %q", disc.Type)
		}
		if err := json.Unmarshal(v, a); err != nil {
			return err
		}
		out[key] = a
	}
	*q = out
	return nil
}

// NoulAnswer is the answer to a NoulQuestion.
type NoulAnswer struct {
	Noul float64 `json:"noul"`
}

func (NoulAnswer) isAnswer() {}

// ChoiceAnswer is the answer to a ChoiceQuestion.
type ChoiceAnswer struct {
	Choice        string             `json:"choice"`
	Probabilities map[string]float64 `json:"probabilities"`
	Confidence    float64            `json:"confidence"`
}

func (ChoiceAnswer) isAnswer() {}

// ScoreAnswer is the answer to a ScoreQuestion.
type ScoreAnswer struct {
	Score         float64            `json:"score"`
	Legend        map[string]string  `json:"legend"`
	Probabilities map[string]float64 `json:"probabilities"`
	Confidence    float64            `json:"confidence"`
}

func (ScoreAnswer) isAnswer() {}
