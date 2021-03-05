package errorx

import (
	"fmt"

	"github.com/pkg/errors"
)

type Error struct {
	err        error
	ApiMessage string
}

func (e Error) String() string {
	return fmt.Sprintf("error: %v, apiMessage:%s", e.err, e.ApiMessage)
}

func Api(err error, message string) Error {
	return Error{
		err:        err,
		ApiMessage: message,
	}
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

func New(message string) error {
	return errors.New(message)
}

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}
