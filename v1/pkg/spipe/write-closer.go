package spipe

import "io"

// SplitWriteCloser defines an io.WriteCloser implementation that writes to and
// can close multiple outputs.
type SplitWriteCloser interface {
	io.WriteCloser

	// IgnoreErrors sets whether or not the split writer should ignore errors
	// returned from secondary writers.
	IgnoreErrors(bool) SplitWriteCloser
}

// NewSplitWriteCloser constructs a new SplitWriteCloser instance with the given
// primary and secondary writers.
func NewSplitWriteCloser(
	raw io.WriteCloser,
	addtl ...io.WriteCloser,
) SplitWriteCloser {
	return &splitWriteCloser{
		primary:   raw,
		secondary: addtl,
	}
}

type splitWriteCloser struct {
	primary    io.WriteCloser
	secondary  []io.WriteCloser
	ignoreErrs bool
}

func (s *splitWriteCloser) Write(p []byte) (n int, err error) {
	if n, err = s.primary.Write(p); err != nil {
		return
	}

	for _, w := range s.secondary {
		if _, err := w.Write(p); err != nil && !s.ignoreErrs {
			return n, err
		}
	}

	return n, nil
}

func (s *splitWriteCloser) Close() (err error) {
	var errs []error

	if e := s.primary.Close(); e != nil {
		errs = append(errs, e)
	}

	for _, w := range s.secondary {
		if e := w.Close(); e != nil && !s.ignoreErrs {
			errs = append(errs, e)
		}
	}

	if len(errs) > 0 {
		err = NewMultiError(errs)
	}

	return
}

func (s *splitWriteCloser) IgnoreErrors(b bool) SplitWriteCloser {
	s.ignoreErrs = b
	return s
}
