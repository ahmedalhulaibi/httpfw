package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type httpError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (he httpError) Error() string {
	return fmt.Sprintf("%d - %s", he.Status, he.Message)
}

func ToHTTPErr(err error) (httpError, bool) {
	e, ok := err.(httpError)
	return e, ok
}

func WriteError(w http.ResponseWriter, err httpError) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	return json.NewEncoder(w).Encode(err)
}

func NotImplementedError() httpError {
	return httpError{http.StatusNotImplemented, "Not implemented"}
}

func InternalServerError(msg string) httpError {
	return httpError{http.StatusInternalServerError, "Internal server error - " + msg}
}

func BadRequest(msg string) httpError {
	return httpError{http.StatusBadRequest, "Bad request - " + msg}
}

func Unauthorized() httpError {
	return httpError{http.StatusUnauthorized, "Unauthorized"}
}

func UnsupportedMediaType(msg string) httpError {
	return httpError{http.StatusUnsupportedMediaType, "Malformed request - " + msg}
}

func RequestEntityTooLarge(msg string) httpError {
	return httpError{http.StatusRequestEntityTooLarge, "Request entity too large - " + msg}
}
