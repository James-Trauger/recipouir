package api

import (
	"encoding/json"
	"testing"
)

func TestJsonErrorEncoding(t *testing.T) {
	msg := "test error"
	expected := "{\"error\":\"" + msg + "\"}"
	err := NewJsonErr(msg, 404)
	raw, marshErr := json.Marshal(&err)
	if marshErr != nil {
		t.Fatal(err)
	}

	if actual := string(raw); actual != expected {
		t.Fatalf("json does not match\nExpected: %s\nReceived: %s\n", expected, actual)
	}
}
