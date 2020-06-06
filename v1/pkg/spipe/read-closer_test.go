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

type testRc struct {
	io.Reader
	cl func() error
}

func (t testRc) Close() error {
	return t.cl()
}

func TestMultiReadCloser_Read(t *testing.T) {
	Convey("MultiReadCloser.Read", t, func() {
		Convey("full read", func() {
			readers := []io.ReadCloser{
				testRc{Reader: strings.NewReader("abc")},
				testRc{Reader: strings.NewReader("def")},
				testRc{Reader: strings.NewReader("ghi")},
				testRc{Reader: strings.NewReader("jkl")},
				testRc{Reader: strings.NewReader("mno")},
			}

			test := spipe.NewMultiReadCloser(readers...)
			buff := make([]byte, 18)

			n, e := test.Read(buff)

			So(e, ShouldBeNil)
			So(n, ShouldEqual, 15)
			So(string(buff[:15]), ShouldResemble, "abcdefghijklmno")

			n, e = test.Read(buff)

			So(e, ShouldEqual, io.EOF)
			So(n, ShouldEqual, 0)
		})

		Convey("chunk read", func() {
			readers := []io.ReadCloser{
				testRc{Reader: strings.NewReader("abc")},
				testRc{Reader: strings.NewReader("def")},
				testRc{Reader: strings.NewReader("ghi")},
				testRc{Reader: strings.NewReader("jkl")},
				testRc{Reader: strings.NewReader("mno")},
			}

			test := spipe.NewMultiReadCloser(readers...)
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

		Convey("aggressive close", func() {
			val := 0
			fun := func() error { val++; return nil }

			readers := []io.ReadCloser{
				testRc{Reader: strings.NewReader("abc"), cl: fun},
				testRc{Reader: strings.NewReader("def"), cl: fun},
				testRc{Reader: strings.NewReader("ghi"), cl: fun},
				testRc{Reader: strings.NewReader("jkl"), cl: fun},
				testRc{Reader: strings.NewReader("mno"), cl: fun},
			}

			test := spipe.NewMultiReadCloser(readers...).CloseImmediately(true)
			buff := make([]byte, 4)

			n, e := test.Read(buff)

			So(e, ShouldBeNil)
			So(n, ShouldEqual, 4)
			So(val, ShouldEqual, 1)

			n, e = test.Read(buff)

			So(e, ShouldBeNil)
			So(n, ShouldEqual, 4)
			So(val, ShouldEqual, 2)

			n, e = test.Read(buff)

			So(e, ShouldBeNil)
			So(n, ShouldEqual, 4)
			So(val, ShouldEqual, 3)

			n, e = test.Read(buff)

			So(e, ShouldBeNil)
			So(n, ShouldEqual, 3)
			So(val, ShouldEqual, 5)

			e = test.Close()

			So(e, ShouldEqual, nil)
			So(val, ShouldEqual, 5)
		})

		Convey("erroring reader", func() {
			okRead := &iotest.ReadCloser{
				ReadableData: []byte("abcdefghi"),
				ReadCounts: []int{3, 3, 3},
			}
			readers := []io.ReadCloser{
				okRead,
				okRead,
				&iotest.ReadCloser{ReadErrors: []error{errors.New("hola")}},
				okRead,
			}

			test := spipe.NewMultiReadCloser(readers...)
			buff := make([]byte, 15)

			_, e := test.Read(buff)

			So(e, ShouldResemble, errors.New("hola"))
		})


		Convey("erroring closer", func() {
			okRead := &iotest.ReadCloser{
				ReadableData: []byte("abcdefghi"),
				ReadCounts: []int{3, 3, 3},
			}

			badClose := &iotest.ReadCloser{
				CloseErrors: []error{errors.New("hola")},
				ReadCounts: []int{0},
			}

			readers := []io.ReadCloser{okRead, okRead, badClose, okRead}

			test := spipe.NewMultiReadCloser(readers...).CloseImmediately(true)
			buff := make([]byte, 15)

			_, e := test.Read(buff)

			So(badClose.CloseCalls, ShouldEqual, 1)
			So(okRead.ReadCalls, ShouldEqual, 2)
			So(okRead.CloseCalls, ShouldEqual, 2)
			So(e, ShouldResemble, errors.New("hola"))
		})
	})
}

func TestMultiReadCloser_Close(t *testing.T) {
	Convey("MultiReadCloser.Close", t, func() {
		Convey("no errors", func() {
			val := 0
			fn := func() error { val++; return nil }

			readers := []io.ReadCloser{
				testRc{cl: fn},
				testRc{cl: fn},
				testRc{cl: fn},
				testRc{cl: fn},
				testRc{cl: fn},
			}

			test := spipe.NewMultiReadCloser(readers...)

			So(test.Close(), ShouldBeNil)
			So(val, ShouldEqual, len(readers))
		})

		Convey("errors", func() {
			fn := func() error { return errors.New("hi") }

			readers := []io.ReadCloser{
				testRc{cl: fn},
				testRc{cl: fn},
				testRc{cl: fn},
				testRc{cl: fn},
				testRc{cl: fn},
			}

			errs, ok := spipe.NewMultiReadCloser(readers...).Close().(spipe.MultiError)

			So(ok, ShouldBeTrue)
			So(errs, ShouldNotBeNil)
			So(errs.Error(), ShouldEqual, "hi\nhi\nhi\nhi\nhi")
			So(errs.Errors(), ShouldResemble, []error{
				errors.New("hi"),
				errors.New("hi"),
				errors.New("hi"),
				errors.New("hi"),
				errors.New("hi"),
			})
		})
	})
}
