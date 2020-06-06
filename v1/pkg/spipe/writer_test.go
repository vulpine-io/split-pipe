package spipe_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/vulpine-io/io-test/v1/pkg/iotest"

	"github.com/vulpine-io/split-pipe/v1/pkg/spipe"
)

func TestSplitWriter_Write(t *testing.T) {
	Convey("SplitWriter.Write", t, func() {

		Convey("happy path", func() {
			a := new(strings.Builder)
			b := new(strings.Builder)
			c := new(strings.Builder)

			test := spipe.NewSplitWriter(a, b, c)
			n, err := test.Write([]byte("hello"))
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 5)
			So(a.String(), ShouldEqual, "hello")
			So(b.String(), ShouldEqual, "hello")
			So(c.String(), ShouldEqual, "hello")
		})

		Convey("failing secondary", func() {

			Convey("without ignore", func() {
				a := new(strings.Builder)
				b := &WriteCloser{WriteErrors: []error{errors.New("hiya!")}}
				c := new(strings.Builder)

				test := spipe.NewSplitWriter(a, b, c)
				n, err := test.Write([]byte("hello"))

				So(err, ShouldResemble, b.WriteErrors[0])
				So(n, ShouldEqual, 5)
				So(a.String(), ShouldEqual, "hello")
				So(b.WrittenBytes, ShouldResemble, []byte("hello"))
				So(c.String(), ShouldEqual, "")
			})

			Convey("with ignore", func() {
				a := new(strings.Builder)
				b := &WriteCloser{WriteErrors: []error{errors.New("hiya!")}}
				c := new(strings.Builder)

				test := spipe.NewSplitWriter(a, b, c).IgnoreErrors(true)
				n, err := test.Write([]byte("hello"))

				So(err, ShouldBeNil)
				So(n, ShouldEqual, 5)
				So(a.String(), ShouldEqual, "hello")
				So(b.WrittenBytes, ShouldResemble, []byte("hello"))
				So(c.String(), ShouldEqual, "hello")
			})
		})

		Convey("failing primary", func() {

			Convey("without ignore", func() {
				a := &WriteCloser{WriteErrors: []error{errors.New("hiya!")}}
				b := new(strings.Builder)
				c := new(strings.Builder)

				test := spipe.NewSplitWriter(a, b, c)
				n, err := test.Write([]byte("hello"))

				So(err, ShouldResemble, a.WriteErrors[0])
				So(n, ShouldEqual, 5)
				So(a.WrittenBytes, ShouldResemble, []byte("hello"))
				So(b.String(), ShouldEqual, "")
				So(c.String(), ShouldEqual, "")
			})

			Convey("with ignore", func() {
				a := &WriteCloser{WriteErrors: []error{errors.New("hiya!")}}
				b := new(strings.Builder)
				c := new(strings.Builder)

				test := spipe.NewSplitWriter(a, b, c)
				n, err := test.Write([]byte("hello"))

				So(err, ShouldResemble, a.WriteErrors[0])
				So(n, ShouldEqual, 5)
				So(a.WrittenBytes, ShouldResemble, []byte("hello"))
				So(b.String(), ShouldEqual, "")
				So(c.String(), ShouldEqual, "")
			})
		})

		Convey("short write", func() {
			Convey("on primary", func() {
				a := &WriteCloser{WriteCounts: []int{1}}
				b := &WriteCloser{}

				test := spipe.NewSplitWriter(a, b)
				n, err := test.Write([]byte("hello"))

				So(err, ShouldEqual, io.ErrShortWrite)
				So(n, ShouldEqual, 1)
				So(b.WriteCalls, ShouldEqual, 0)
			})

			Convey("on secondary", func() {
				Convey("without ignore errors", func() {
					a := &WriteCloser{}
					b := &WriteCloser{WriteCounts: []int{1}}

					test := spipe.NewSplitWriter(a, b)
					n, err := test.Write([]byte("hello"))

					So(string(a.WrittenBytes), ShouldEqual, "hello")
					So(n, ShouldEqual, 1)
					So(err, ShouldEqual, io.ErrShortWrite)
				})

				Convey("with ignore errors", func() {
					a := &WriteCloser{}
					b := &WriteCloser{WriteCounts: []int{1}}

					test := spipe.NewSplitWriter(a, b).IgnoreErrors(true)
					n, err := test.Write([]byte("hello"))

					So(string(a.WrittenBytes), ShouldEqual, "hello")
					So(n, ShouldEqual, 5)
					So(err, ShouldBeNil)
				})

			})
		})

	})
}
