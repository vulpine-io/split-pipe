package spipe_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vulpine-io/split-pipe/v1/pkg/spipe"
	"strings"
	"testing"
)

func TestSplitWriter_Write(t *testing.T) {
	Convey("SplitWriter.Write", t, func() {
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
}
