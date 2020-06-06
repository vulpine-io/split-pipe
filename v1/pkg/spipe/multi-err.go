package spipe

import "strings"

// MultiError wraps a slice of errors into a single error type.
type MultiError interface {
	error

	// Errors returns the original errors backing this error.
	Errors() []error
}

// NewMultiError constructs a new MultiError instance from the given slice of
// errors.
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
