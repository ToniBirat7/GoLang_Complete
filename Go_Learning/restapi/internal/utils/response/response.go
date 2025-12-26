package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

type Response struct {
	Status string
	Error  string
}

const (
	StatusOk    = "OK"
	StatusError = "ERROR"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMessages []string

	for _, e := range errs {
		switch e.ActualTag() {
		case "required":
			errMessages = append(errMessages, fmt.Sprintf("field %s is required", e.Field()))
		default:
			errMessages = append(errMessages, fmt.Sprintf("field %s failed on %s", e.Field(), e.ActualTag()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMessages, "; "),
	}
}
