package spipe_test

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vulpine-io/split-pipe/v1/pkg/spipe"
	"strings"
	"testing"
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
