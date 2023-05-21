package brutil

import (
	"encoding/json"
	"net/http"
)

func SendJSON(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return Send(w, statusCode, data)
}

func Send(w http.ResponseWriter, statusCode int, data []byte) error {
	w.WriteHeader(statusCode)
	_, err := w.Write(data)
	return err
}
