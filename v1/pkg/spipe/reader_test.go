package spipe_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vulpine-io/split-pipe/v1/pkg/spipe"
	"io"
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
