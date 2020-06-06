package spipe

import "io"

// MultiReadCloser defines an io.ReadCloser implementation that can read from
// and close multiple input streams as if they were one long stream.
//
// The input streams will be read to completion in the order they are given to
// the MultiReader instance.
//
// If the CloseImmediately option is set to true, the given streams will be
// closed, in order, as soon as they hit EOF.
type MultiReadCloser interface {
	io.ReadCloser

	// CloseImmediately controls whether the input readers will be closed as soon
	// as they are consumed rather than waiting for a Close call.
	CloseImmediately(bool) MultiReadCloser
}

// NewMultiReadCloser returns a new MultiReadCloser instance that will read from
// the given inputs in the order they are passed.
func NewMultiReadCloser(inputs ...io.ReadCloser) MultiReadCloser {
	return &multiReadCloser{inputs: inputs}
}

type multiReadCloser struct {
	inputs   []io.ReadCloser
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

// Read attempts to fill the given buffer by reading from one or more available
// streams until it runs out of input, or the len(p) bytes have been read.
//
// If read returns a bytes-read count > 0, it will not return an EOF.  If a read
// reaches the end of all available inputs, the total number of bytes read will
// be returned with a nil error, and subsequent calls will return an EOF.
//
// This method behaves differently than the current version of `io.MultiWriter`
// (as of v1.14.4) in that the `io.MultiWriter` implementation will read less
// than len(p) bytes any time it encounters the end of a single stream, whereas
// this method will automatically continue on to the next stream in a single
// call to Read in order to fill the input buffer.
//
// In practice, the following code
//     buffer := make([]byte, 512)
//     spipe.NewMultiReader(reader1, reader2).Read(buffer)
//
// is functionally closer to
//     buffer := bytes.NewBuffer(make([]byte, 0, 512))
//     io.Copy(buffer, io.MultiReader(reader1, reader2))
//
// than it is to
//     buffer := make([]byte, 512)
//     io.MultiReader(reader1, reader2).Read(buffer)
func (m *multiReadCloser) Read(p []byte) (totalRead int, err error) {
	return internalRead(m, p)
}

func (m *multiReadCloser) hasNext() bool {
	return len(m.inputs) > 0
}

func (m *multiReadCloser) nextInput() io.Reader {
	return m.inputs[0]
}

func (m *multiReadCloser) popInput() (err error) {
	// 0 case is not possible due to hasNext call in read
	if len(m.inputs) == 1 {
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
