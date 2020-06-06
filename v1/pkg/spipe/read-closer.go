package spipe

import "io"

type MultiReadCloser interface {
	io.ReadCloser
}

func NewMultiReadCloser(inputs ...io.ReadCloser) MultiReadCloser {
	return &multiReadCloser{inputs}
}

type multiReadCloser struct {
	inputs []io.ReadCloser
}

func (m *multiReadCloser) Close() (err error) {
	var errs []error

	for _, r := range m.inputs {
		if e := r.Close(); e != nil {
			errs = append(errs, e)
		}
	}

	if len(errs) > 0 {
		err = NewMultiError(errs)
	}

	return
}

func (m *multiReadCloser) Read(p []byte) (n int, err error) {
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

func (m *multiReadCloser) hasNext() bool {
	return len(m.inputs) > 0
}

func (m *multiReadCloser) nextInput() io.Reader {
	return m.inputs[0]
}

func (m *multiReadCloser) popInput() {
	if len(m.inputs) < 2 {
		m.inputs = nil
		return
	}

	m.inputs = m.inputs[1:]
}

