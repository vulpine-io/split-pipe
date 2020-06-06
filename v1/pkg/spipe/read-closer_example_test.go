package spipe_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/vulpine-io/split-pipe/v1/pkg/spipe"
)

type PrintCloser struct {
	io.Reader
	index int
}

func (p PrintCloser) Close() error {
	fmt.Printf("Closed %d!\n", p.index)
	return nil
}

func ExampleMultiReadCloser_Read() {
	input1 := PrintCloser{strings.NewReader("hello"), 1}
	input2 := PrintCloser{bytes.NewReader([]byte{' '}), 2}
	input3 := PrintCloser{strings.NewReader("world"), 3}

	buffer := make([]byte, 11)
	reader := spipe.NewMultiReadCloser(input1, input2, input3).
		CloseImmediately(true)

	fmt.Println(reader.Read(buffer))
	fmt.Println(string(buffer))

	reader.Close()

	// Output: Closed 1!
	// Closed 2!
	// 11 <nil>
	// hello world
	// Closed 3!
}
