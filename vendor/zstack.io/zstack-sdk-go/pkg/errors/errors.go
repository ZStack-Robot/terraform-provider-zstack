package errors

import "github.com/pkg/errors"

type Error string

func (e Error) Error() string {
	return string(e)
}

func Wrap(err error, msg string) error {
	return errors.Wrap(err, msg)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

func Cause(err error) error {
	return errors.Cause(err)
}

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}
