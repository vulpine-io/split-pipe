package spipe_test

import (
	"errors"
	"io"
	"strings"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/vulpine-io/io-test/v1/pkg/iotest"
)

func tReaderComm(construct func(interface{}) io.Reader) {
	Convey("full read", func() {
		readers := []io.Reader{
			strings.NewReader("abc"),
			strings.NewReader("def"),
			strings.NewReader("ghi"),
			strings.NewReader("jkl"),
			strings.NewReader("mno"),
		}

		test := construct(readers)
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

		test := construct(readers)
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

		test := construct(readers)
		buff := make([]byte, 15)

		_, e := test.Read(buff)

		So(e, ShouldResemble, errors.New("hola"))
	})

	Convey("repeating reader", func() {
		readers := []io.Reader{
			io.MultiReader(
				strings.NewReader("abc"),
				strings.NewReader("def")),
			strings.NewReader("ghi"),
		}

		test := construct(readers)
		buff := make([]byte, 9)

		n, err := test.Read(buff)

		So(err, ShouldBeNil)
		So(n, ShouldEqual, 9)
		So(string(buff), ShouldEqual, "abcdefghi")
	})
}
