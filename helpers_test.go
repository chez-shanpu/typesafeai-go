// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of typesafeai-go

package typesafeai_test

import (
	"encoding/json/v2"
	"reflect"
	"testing"
)

func assertJSONEqual(t *testing.T, give []byte, want string) {
	t.Helper()
	var gotValue, wantValue any
	if err := json.Unmarshal(give, &gotValue); err != nil {
		t.Fatalf("decode actual JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatalf("decode expected JSON: %v", err)
	}
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Errorf("JSON = %s, want %s", give, want)
	}
}
