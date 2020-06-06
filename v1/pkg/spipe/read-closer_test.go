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
		fun := func(i interface{}) io.Reader {
			tmp := i.([]io.Reader)
			par := make([]io.ReadCloser, len(tmp))
			for i, r := range tmp {
				par[i] = testRc{Reader: r}
			}
			return spipe.NewMultiReadCloser(par...)
		}

		tReaderComm(fun)

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

		Convey("erroring closer", func() {
			okRead1 := &iotest.ReadCloser{
				ReadableData: []byte("abcdefghi"),
				ReadCounts:   []int{3, 0},
				ReadErrors:   []error{nil, io.EOF},
			}

			okRead2 := &iotest.ReadCloser{
				ReadableData: []byte("abcdefghi"),
				ReadCounts:   []int{3, 0},
				ReadErrors:   []error{nil, io.EOF},
			}

			badClose := &iotest.ReadCloser{
				ReadErrors:  []error{io.EOF},
				CloseErrors: []error{errors.New("hola")},
				ReadCounts:  []int{0},
			}

			readers := []io.ReadCloser{okRead1, okRead2, badClose, okRead1}

			test := spipe.NewMultiReadCloser(readers...).CloseImmediately(true)
			buff := make([]byte, 15)

			_, e := test.Read(buff)

			So(okRead1.ReadCalls, ShouldEqual, 2)
			So(okRead1.CloseCalls, ShouldEqual, 1)
			So(okRead2.ReadCalls, ShouldEqual, 2)
			So(okRead2.CloseCalls, ShouldEqual, 1)
			So(badClose.ReadCalls, ShouldEqual, 1)
			So(badClose.CloseCalls, ShouldEqual, 1)
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
