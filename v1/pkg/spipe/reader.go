package spipe

import "io"

// MultiReader defines an io.Reader implementation that can read from multiple
// input streams as if they were one long stream.
//
// The input streams will be read to completion in the order they are given to
// the MultiReader instance.
type MultiReader interface {
	io.Reader
}

// NewMultiReader returns a new MultiReader instance that will read from the
// given inputs in the order they are passed.
//
// This operates in a manner different to `io.MultiReader` in that this
// implementation will proactively start reading from a secondary stream in a
// single Read call.  The stdlib MultiReader returns from the Read call as soon
// as a single, non-empty stream is exhausted.
func NewMultiReader(inputs ...io.Reader) MultiReader {
	return &multiReader{inputs}
}

type multiReader struct {
	inputs []io.Reader
}

func (m *multiReader) Read(p []byte) (n int, err error) {
	if !m.hasNext() {
		return 0, io.EOF
	}

	ln := len(p)
	n, err = m.nextInput().Read(p)
	if err != nil && err != io.EOF {
		return n, err
	}

	if n < ln {
		m.popInput()

		m, err := m.Read(p[n:])

		if n > 0 && err == io.EOF {
			err = nil
		}

		n += m

		return n, err
	}

	return
}

func (m *multiReader) hasNext() bool {
	return len(m.inputs) > 0
}

func (m *multiReader) nextInput() io.Reader {
	return m.inputs[0]
}

func (m *multiReader) popInput() {
	if len(m.inputs) < 2 {
		m.inputs = nil
		return
	}

	// explicitly free the reference to the exhausted reader.
	m.inputs[0] = nil

	// subset the input readers
	m.inputs = m.inputs[1:]
}
