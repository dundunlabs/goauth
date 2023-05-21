package brutil

import (
	"bytes"
	"io"
	"testing"
)

type mockReader struct{}

func (r *mockReader) Read(p []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestNewBodyParserError(t *testing.T) {
	r := new(mockReader)
	want := io.ErrUnexpectedEOF
	if _, err := NewBodyParser(r); err != want {
		t.Errorf("expected error: %v, got: %v", want, err)
	}
}

func TestParseJSON(t *testing.T) {
	parser, err := NewBodyParser(bytes.NewBufferString(`{"foo":"bar"}`))
	if err != nil {
		t.Error(err)
		return
	}

	type Body struct {
		Foo string `validate:"required"`
	}
	var body Body
	if err := parser.ParseJSON(&body); err != nil {
		t.Error(err)
		return
	}

	if got, want := body.Foo, "bar"; got != want {
		t.Errorf("wrong result. expected: %s, got: %s", want, got)
		return
	}
}

func TestParseJSONError(t *testing.T) {
	parser, err := NewBodyParser(bytes.NewBufferString(`{"foo":bar}`))
	if err != nil {
		t.Error(err)
		return
	}

	var body map[string]any
	if err := parser.ParseJSON(&body); err == nil {
		t.Error("parsing JSON should be failed")
		return
	}
}
