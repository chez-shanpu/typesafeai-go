package typesafeai_test

import (
	"encoding/json/v2"
	"testing"

	"github.com/chez-shanpu/typesafeai-go"
)

func TestQuestionMarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		give typesafeai.Question
		want string
	}{
		{
			name: "noul without criteria",
			give: &typesafeai.NoulQuestion{Instructions: "Is this relevant?"},
			want: `{"type":"noul","instructions":"Is this relevant?"}`,
		},
		{
			name: "noul with criteria",
			give: &typesafeai.NoulQuestion{
				Instructions: "Is this relevant?",
				Criteria:     &typesafeai.NoulCriteria{True: new("relevant"), False: new("irrelevant")},
			},
			want: `{"type":"noul","instructions":"Is this relevant?","criteria":{"true":"relevant","false":"irrelevant"}}`,
		},
		{
			name: "noul with only true criterion",
			give: &typesafeai.NoulQuestion{
				Instructions: "Evaluate",
				Criteria:     &typesafeai.NoulCriteria{True: new("relevant")},
			},
			want: `{"type":"noul","instructions":"Evaluate","criteria":{"true":"relevant"}}`,
		},
		{
			name: "choice preserves null and empty description",
			give: &typesafeai.ChoiceQuestion{
				Instructions: map[string]any{"task": "Choose", "labels": []string{"a", "b", "c"}},
				Criteria:     map[string]*string{"a": nil, "b": new(""), "c": new("third")},
			},
			want: `{"type":"choice","instructions":{"task":"Choose","labels":["a","b","c"]},"criteria":{"a":null,"b":"","c":"third"}}`,
		},
		{
			name: "score preserves order and array instructions",
			give: &typesafeai.ScoreQuestion{
				Instructions: []any{"Evaluate", map[string]string{"axis": "quality"}},
				Criteria:     []string{"poor", "fair", "good"},
			},
			want: `{"type":"score","instructions":["Evaluate",{"axis":"quality"}],"criteria":["poor","fair","good"]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.give)
			if err != nil {
				t.Fatalf("Marshal() error: %v", err)
			}
			assertJSONEqual(t, got, tt.want)
		})
	}
}
