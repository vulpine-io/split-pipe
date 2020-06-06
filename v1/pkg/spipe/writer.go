package spipe

import "io"

// SplitWriter defines an io.Writer implementation that writes to multiple
// outputs.
//
// SplitWriter's implementation differs from `io.MultiWriter` in that it
// provides the option to ignore errors from secondary writers which lessens the
// need for composed wrappers for that particular use case.  If you do not wish
// to ignore secondary writer errors, then SplitWriter is effectively the same
// as `io.MultiWriter`.
type SplitWriter interface {
	io.Writer

	// IgnoreErrors sets whether or not the split writer should ignore errors
	// returned from secondary writers.
	IgnoreErrors(bool) SplitWriter
}

// NewSplitWriter constructs a new SplitWriter instance with the given primary
// and secondary writers.
func NewSplitWriter(raw io.Writer, addtl ...io.Writer) SplitWriter {
	return &splitWriter{primary: raw, secondary: addtl}
}

type splitWriter struct {
	primary    io.Writer
	secondary  []io.Writer
	ignoreErrs bool
}

func (s *splitWriter) Write(p []byte) (n int, err error) {
	if n, err = s.primary.Write(p); err != nil {
		return
	}

	if n < len(p) {
		err = io.ErrShortWrite
		return
	}

	for _, w := range s.secondary {
		m, err := w.Write(p)

		if err != nil && !s.ignoreErrs {
			return n, err
		}

		if m < len(p) && !s.ignoreErrs {
			return m, io.ErrShortWrite
		}
	}

	return n, nil
}

func (s *splitWriter) IgnoreErrors(b bool) SplitWriter {
	s.ignoreErrs = b
	return s
}
