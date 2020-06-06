package spipe

import "io"

// SplitWriter defines an io.Writer implementation that writes to multiple
// outputs.
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

	for _, w := range s.secondary {
		if _, err := w.Write(p); err != nil && !s.ignoreErrs {
			return n, err
		}
	}

	return n, nil
}

func (s *splitWriter) IgnoreErrors(b bool) SplitWriter {
	s.ignoreErrs = b
	return s
}
