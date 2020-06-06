package spipe_test

import (
	"bytes"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vulpine-io/split-pipe/v1/pkg/spipe"
	"io"
	"strings"
	"testing"
)

func TestMultiReader_Read(t *testing.T) {
	Convey("MultiReader.Read", t, func() {
		fun := func(i interface{}) io.Reader {
			return spipe.NewMultiReader(i.([]io.Reader)...)
		}
		tReaderComm(fun)
	})
}

func ExampleMultiReader_Read() {
	input1 := strings.NewReader("hello")
	input2 := bytes.NewReader([]byte{' '})
	input3 := strings.NewReader("world")

	buffer := make([]byte, 11)
	reader := spipe.NewMultiReader(input1, input2, input3)

	fmt.Println(reader.Read(buffer))
	fmt.Println(string(buffer))

	// Output: 11 <nil>
	// hello world
}
