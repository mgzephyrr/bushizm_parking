package pgstorage

import "fmt"

type ErrorDoesNotExist struct {
	What  string
	Inner error
}

func (e ErrorDoesNotExist) Error() string {
	return fmt.Sprintf("%s does not exist", e.What)
}

func (e ErrorDoesNotExist) Unwrap() error {
	return e.Inner
}
