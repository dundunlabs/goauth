package brutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSend(t *testing.T) {
	w := httptest.NewRecorder()
	if err := Send(w, http.StatusCreated, nil); err != nil {
		t.Error(err)
		return
	}
	if got, want := w.Result().StatusCode, http.StatusCreated; got != want {
		t.Errorf("wrong status. expected: %d, got: %d", got, want)
		return
	}
}

func TestSendJSON(t *testing.T) {
	w := httptest.NewRecorder()
	if err := SendJSON(w, http.StatusOK, map[string]string{"foo": "bar"}); err != nil {
		t.Error(err)
		return
	}
	if got, want := w.Result().StatusCode, http.StatusOK; got != want {
		t.Errorf("wrong status. expected: %d, got: %d", got, want)
		return
	}
	if got, want := w.HeaderMap.Get("Content-Type"), "application/json"; got != want {
		t.Errorf("wrong content-type. expected: %s, got: %s", got, want)
		return
	}
	if got, want := w.Body.String(), `{"foo":"bar"}`; got != want {
		t.Errorf("wrong body. expected: %s, got: %s", got, want)
		return
	}
}

func TestSendJSONError(t *testing.T) {
	w := httptest.NewRecorder()
	err := SendJSON(w, http.StatusOK, make(chan int))
	switch err := err.(type) {
	case *json.UnsupportedTypeError:
		return
	default:
		t.Errorf("expected error: %T, got: %T", json.UnsupportedTypeError{}, err)
		return
	}
}
