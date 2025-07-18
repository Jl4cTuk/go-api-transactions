package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(
				errMsgs,
				fmt.Sprintf("field %s is a required field", err.Field()),
			)
		case "gt":
			errMsgs = append(
				errMsgs,
				fmt.Sprintf("field %s must be greater than 0", err.Field()),
			)
		case "nefield":
			errMsgs = append(
				errMsgs,
				"cant send to same wallet",
			)
		default:
			errMsgs = append(
				errMsgs,
				fmt.Sprintf("field %s is not valid", err.Field()),
			)

		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
