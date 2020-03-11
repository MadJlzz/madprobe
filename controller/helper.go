package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	maxBodyBytes = 1048576
	jsonMimeType = "application/json"
)

type malformedContent struct {
	status int
	msg    string
}

func (mr *malformedContent) Error() string {
	return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, req *http.Request, dst interface{}) error {
	if req.Header.Get("Content-Type") != jsonMimeType {
		msg := "Content-Type header is not application/json"
		return &malformedContent{
			status: http.StatusUnsupportedMediaType,
			msg:    msg,
		}
	}

	req.Body = http.MaxBytesReader(w, req.Body, maxBodyBytes)

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedContent{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedContent{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedContent{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedContent{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedContent{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedContent{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	if dec.More() {
		msg := "Request body must only contain a single JSON object"
		return &malformedContent{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}

func encodeJSONBody(w http.ResponseWriter, src interface{}) error {
	w.Header().Set("Content-Type", jsonMimeType)

	enc := json.NewEncoder(w)
	err := enc.Encode(&src)
	if err != nil {
		var syntaxError *json.SyntaxError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Model contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedContent{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Model body contains badly-formed JSON")
			return &malformedContent{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Model body must not be empty"
			return &malformedContent{status: http.StatusBadRequest, msg: msg}

		default:
			return err
		}
	}
	return nil
}
