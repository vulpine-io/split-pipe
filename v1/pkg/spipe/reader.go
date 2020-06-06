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
func NewMultiReader(inputs ...io.Reader) MultiReader {
	return &multiReader{inputs}
}

type multiReader struct {
	inputs []io.Reader
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
func (m *multiReader) Read(p []byte) (totalRead int, err error) {
	return internalRead(m, p)
}

func (m *multiReader) hasNext() bool {
	return len(m.inputs) > 0
}

func (m *multiReader) nextInput() io.Reader {
	return m.inputs[0]
}

func (m *multiReader) popInput() (_ error) {
	if len(m.inputs) < 2 {
		m.inputs = nil
		return
	}

	// explicitly free the reference to the exhausted reader.
	m.inputs[0] = nil

	// subset the input readers
	m.inputs = m.inputs[1:]

	return
}
