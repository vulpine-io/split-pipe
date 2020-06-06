package spipe

import "io"

type MultiReader interface {
	io.Reader
}

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

	ln  := len(p)
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

	m.inputs = m.inputs[1:]
}

