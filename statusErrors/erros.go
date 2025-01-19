package statuserrors

import (
	"fmt"
	"net/http"
)

type StatusError interface {
	error
	Status() int
}

type statusError int

var (
	ErrBadRequest          statusError = http.StatusBadRequest
	ErrNotFound            statusError = http.StatusNotFound
	ErrInternalServerError statusError = http.StatusInternalServerError
)

func (se statusError) Error() string {
	return http.StatusText(int(se))
}

func (se statusError) Status() int {
	return int(se)
}

func NewStatusError(status int, msg string) error {
	return fmt.Errorf("%w: %s", statusError(status), msg)
}

func NewNotFoundError(msg string) error {
	return NewStatusError(http.StatusBadRequest, msg)
}

func NewBadRequestError(msg string) error {
	return NewStatusError(http.StatusBadRequest, msg)
}

func NewInternalServerError(msg string) error {
	return NewStatusError(http.StatusInternalServerError, msg)
}
