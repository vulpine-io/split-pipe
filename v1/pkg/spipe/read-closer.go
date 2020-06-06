package spipe

import "io"

type MultiReadCloser interface {
	io.ReadCloser

	// CloseImmediately controls whether the input readers will be closed as soon
	// as they are consumed rather than waiting for a Close call.
	CloseImmediately(bool) MultiReadCloser
}

func NewMultiReadCloser(inputs ...io.ReadCloser) MultiReadCloser {
	return &multiReadCloser{inputs: inputs}
}

type multiReadCloser struct {
	inputs []io.ReadCloser
	aggClose bool
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

func (m *multiReadCloser) CloseImmediately(b bool) MultiReadCloser {
	m.aggClose = b
	return m
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
		if err := m.popInput(); err != nil {
			return n, err
		}

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

func (m *multiReadCloser) popInput() (err error) {
	switch len(m.inputs) {
	case 0:
		return
	case 1:
		if m.aggClose {
			err = m.inputs[0].Close()
		}
		m.inputs = nil
		return
	}

	if m.aggClose {
		err = m.inputs[0].Close()
	}

	m.inputs = m.inputs[1:]

	return
}

