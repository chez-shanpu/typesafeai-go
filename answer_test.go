package typesafeai_test

import (
	"encoding/json/v2"
	"reflect"
	"testing"

	"github.com/chez-shanpu/typesafeai-go"
)

func TestAnswerUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		give string
		want typesafeai.Answer
	}{
		{
			name: "noul",
			give: `{"type":"noul","noul":0.8}`,
			want: &typesafeai.NoulAnswer{Noul: 0.8},
		},
		{
			name: "choice",
			give: `{"type":"choice","choice":"a","probabilities":{"a":0.75,"b":0.25},"confidence":0.9}`,
			want: &typesafeai.ChoiceAnswer{
				Choice: "a", Probabilities: map[string]float64{"a": 0.75, "b": 0.25}, Confidence: 0.9,
			},
		},
		{
			name: "score",
			give: `{"type":"score","score":0.75,"legend":{"0":"poor","1":"good"},"probabilities":{"0":0.25,"1":0.75},"confidence":0.9}`,
			want: &typesafeai.ScoreAnswer{
				Score: 0.75, Legend: map[string]string{"0": "poor", "1": "good"},
				Probabilities: map[string]float64{"0": 0.25, "1": 0.75}, Confidence: 0.9,
			},
		},
		{
			name: "zero noul",
			give: `{"type":"noul","noul":0}`,
			want: &typesafeai.NoulAnswer{},
		},
		{
			name: "zero choice confidence and probability",
			give: `{"type":"choice","choice":"a","probabilities":{"a":0},"confidence":0}`,
			want: &typesafeai.ChoiceAnswer{Choice: "a", Probabilities: map[string]float64{"a": 0}},
		},
		{
			name: "zero score confidence and probability",
			give: `{"type":"score","score":0,"legend":{"0":"poor"},"probabilities":{"0":0},"confidence":0}`,
			want: &typesafeai.ScoreAnswer{
				Legend: map[string]string{"0": "poor"}, Probabilities: map[string]float64{"0": 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got typesafeai.QuestionKeyToAnswer
			if err := json.Unmarshal([]byte(`{"question":`+tt.give+`}`), &got); err != nil {
				t.Fatalf("Unmarshal() error: %v", err)
			}
			want := typesafeai.QuestionKeyToAnswer{"question": tt.want}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("answers = %#v, want %#v", got, want)
			}
		})
	}
}

func TestAnswerUnmarshalJSONMalformed(t *testing.T) {
	var got typesafeai.QuestionKeyToAnswer
	if err := json.Unmarshal([]byte(`{"question":`), &got); err == nil {
		t.Fatal("Unmarshal() succeeded for malformed JSON")
	}
}
