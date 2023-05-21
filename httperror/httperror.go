package httperror

import (
	"net/http"

	"github.com/dundunlabs/goauth/brutil"
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bunrouter"
)

func ErrorMiddleware(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		if err := next(w, req); err != nil {
			httperr := handlerError(err)
			brutil.SendJSON(w, httperr.statusCode, httperr)

			if httperr.statusCode == http.StatusInternalServerError {
				return err
			}
		}

		return nil
	}
}

func handlerError(err error) HTTPError {
	switch err := err.(type) {
	case HTTPError:
		return err
	case validator.ValidationErrors:
		return HTTPError{
			statusCode: http.StatusBadRequest,
			Code:       "bad_request",
			Message:    err.Error(),
		}
	default:
		return HTTPError{
			statusCode: http.StatusInternalServerError,
			Code:       "internal_server_error",
			Message:    err.Error(),
		}
	}
}

func New(statusCode int, code string, message string) HTTPError {
	return HTTPError{
		statusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

type HTTPError struct {
	statusCode int

	Code    string `json:"code"`
	Message string `json:"message"`
}

func (err HTTPError) Error() string {
	return err.Message
}
