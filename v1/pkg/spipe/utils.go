package spipe

import "io"

// NewEOFReader returns an instance of io.Reader that returns an io.EOF error on
// first read.
func NewEOFReader() io.Reader {
	return eofer{}
}

// NewEOFReadCloser returns an instance of io.ReadCloser that returns an io.EOF
// error on first read.
func NewEOFReadCloser() io.ReadCloser {
	return eofer{}
}

type eofer struct {}

func (e eofer) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (e eofer) Close() error {
	return nil
}
