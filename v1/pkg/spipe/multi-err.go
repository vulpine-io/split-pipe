package spipe

import "strings"

type MultiError interface {
	error

	Errors() []error
}

func NewMultiError(errs []error) MultiError {
	return &multiError{errs}
}

type multiError struct {
	errs []error
}

func (m *multiError) Error() string {
	out := strings.Builder{}

	for i, e := range m.errs {
		if i > 0 {
			out.WriteByte('\n')
		}

		out.WriteString(e.Error())
	}

	return out.String()
}

func (m *multiError) Errors() []error {
	return m.errs
}

