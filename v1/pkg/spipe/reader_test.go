package spipe_test

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vulpine-io/io-test/v1/pkg/iotest"
	"github.com/vulpine-io/split-pipe/v1/pkg/spipe"
	"io"
	"strings"
	"testing"
)

func TestMultiReader_Read(t *testing.T) {
	Convey("MultiReader.Read", t, func() {
		Convey("full read", func() {
			readers := []io.Reader{
				strings.NewReader("abc"),
				strings.NewReader("def"),
				strings.NewReader("ghi"),
				strings.NewReader("jkl"),
				strings.NewReader("mno"),
			}

			test := spipe.NewMultiReader(readers...)
			buff := make([]byte, 15)

			n, e := test.Read(buff)

			So(e, ShouldBeNil)
			So(n, ShouldEqual, 15)
			So(string(buff), ShouldEqual, "abcdefghijklmno")

			n, e = test.Read(buff)

			So(e, ShouldEqual, io.EOF)
			So(n, ShouldEqual, 0)
		})

		Convey("chunk read", func() {
			readers := []io.Reader{
				strings.NewReader("abc"),
				strings.NewReader("def"),
				strings.NewReader("ghi"),
				strings.NewReader("jkl"),
				strings.NewReader("mno"),
			}

			test := spipe.NewMultiReader(readers...)
			buff := make([]byte, 4)

			n, e := test.Read(buff)

			So(e, ShouldBeNil)
			So(n, ShouldEqual, 4)
			So(string(buff), ShouldEqual, "abcd")

			n, e = test.Read(buff)

			So(e, ShouldBeNil)
			So(n, ShouldEqual, 4)
			So(string(buff), ShouldEqual, "efgh")

			n, e = test.Read(buff)

			So(e, ShouldBeNil)
			So(n, ShouldEqual, 4)
			So(string(buff), ShouldEqual, "ijkl")

			n, e = test.Read(buff)

			So(e, ShouldBeNil)
			So(n, ShouldEqual, 3)
			So(string(buff), ShouldEqual, "mnol")

			n, e = test.Read(buff)

			So(e, ShouldEqual, io.EOF)
			So(n, ShouldEqual, 0)
		})

		Convey("erroring reader", func() {
			readers := []io.Reader{
				strings.NewReader("abc"),
				strings.NewReader("def"),
				&iotest.ReadCloser{ReadErrors: []error{errors.New("hola")}},
				strings.NewReader("ghi"),
			}

			test := spipe.NewMultiReader(readers...)
			buff := make([]byte, 15)

			_, e := test.Read(buff)

			So(e, ShouldResemble, errors.New("hola"))
		})
	})
}
