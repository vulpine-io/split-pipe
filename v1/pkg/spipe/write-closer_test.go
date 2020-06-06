package spipe_test

import (
	"errors"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/vulpine-io/io-test/v1/pkg/iotest"

	"github.com/vulpine-io/split-pipe/v1/pkg/spipe"
)

type testWC struct {
	strings.Builder
	cl func() error
}

func (t *testWC) Close() error {
	return t.cl()
}

func TestSplitWriteCloser_Write(t *testing.T) {
	Convey("SplitWriteCloser.Write", t, func() {
		a := new(testWC)
		b := new(testWC)
		c := new(testWC)

		test := spipe.NewSplitWriteCloser(a, b, c)
		n, err := test.Write([]byte("hello"))
		So(err, ShouldBeNil)
		So(n, ShouldEqual, 5)
		So(a.String(), ShouldEqual, "hello")
		So(b.String(), ShouldEqual, "hello")
		So(c.String(), ShouldEqual, "hello")

		Convey("failing secondary", func() {

			Convey("without ignore", func() {
				a := new(WriteCloser)
				b := &WriteCloser{WriteErrors: []error{errors.New("hiya!")}}
				c := new(WriteCloser)

				test := spipe.NewSplitWriteCloser(a, b, c)
				n, err := test.Write([]byte("hello"))

				So(err, ShouldResemble, b.WriteErrors[0])
				So(n, ShouldEqual, 5)
				So(a.WrittenBytes, ShouldResemble, []byte("hello"))
				So(b.WrittenBytes, ShouldResemble, []byte("hello"))
				So(c.WrittenBytes, ShouldBeEmpty)
			})

			Convey("with ignore", func() {
				a := new(WriteCloser)
				b := &WriteCloser{WriteErrors: []error{errors.New("hiya!")}}
				c := new(WriteCloser)

				test := spipe.NewSplitWriteCloser(a, b, c).IgnoreErrors(true)
				n, err := test.Write([]byte("hello"))

				So(err, ShouldBeNil)
				So(n, ShouldEqual, 5)
				So(a.WrittenBytes, ShouldResemble, []byte("hello"))
				So(b.WrittenBytes, ShouldResemble, []byte("hello"))
				So(c.WrittenBytes, ShouldResemble, []byte("hello"))
			})
		})

		Convey("failing primary", func() {

			Convey("without ignore", func() {
				a := &WriteCloser{WriteErrors: []error{errors.New("hiya!")}}
				b := new(WriteCloser)
				c := new(WriteCloser)

				test := spipe.NewSplitWriteCloser(a, b, c)
				n, err := test.Write([]byte("hello"))

				So(err, ShouldResemble, a.WriteErrors[0])
				So(n, ShouldEqual, 5)
				So(a.WrittenBytes, ShouldResemble, []byte("hello"))
				So(b.WrittenBytes, ShouldBeEmpty)
				So(c.WrittenBytes, ShouldBeEmpty)
			})

			Convey("with ignore", func() {
				a := &WriteCloser{WriteErrors: []error{errors.New("hiya!")}}
				b := new(WriteCloser)
				c := new(WriteCloser)

				test := spipe.NewSplitWriteCloser(a, b, c)
				n, err := test.Write([]byte("hello"))

				So(err, ShouldResemble, a.WriteErrors[0])
				So(n, ShouldEqual, 5)
				So(a.WrittenBytes, ShouldResemble, []byte("hello"))
				So(b.WrittenBytes, ShouldBeEmpty)
				So(c.WrittenBytes, ShouldBeEmpty)
			})
		})
	})
}

func TestSplitWriteCloser_Close(t *testing.T) {
	Convey("SplitWriteCloser.Write", t, func() {
		Convey("no errors", func() {
			val := 0
			fun := func() error { val++; return nil }
			a := &testWC{cl: fun}
			b := &testWC{cl: fun}
			c := &testWC{cl: fun}

			test := spipe.NewSplitWriteCloser(a, b, c)
			err := test.Close()

			So(err, ShouldBeNil)
			So(val, ShouldEqual, 3)
		})

		Convey("errors", func() {
			fun := func() error { return errors.New("hi") }
			a := &testWC{cl: fun}
			b := &testWC{cl: fun}
			c := &testWC{cl: fun}

			err, ok := spipe.NewSplitWriteCloser(a, b, c).Close().(spipe.MultiError)

			So(ok, ShouldBeTrue)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "hi\nhi\nhi")
		})
	})
}
