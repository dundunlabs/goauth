package brutil

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

func NewBodyParser(r io.Reader) (*BodyParser, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return &BodyParser{
		data:     data,
		validate: validator.New(),
	}, nil
}

type BodyParser struct {
	data     []byte
	validate *validator.Validate
}

func (bp *BodyParser) ParseJSON(v any) error {
	if err := json.Unmarshal(bp.data, v); err != nil {
		return err
	}
	return bp.validate.Struct(v)
}
