// History: Apr 26 14 tcolar Creation

package album

import (
	"log"
	"os"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tcolar/utils"
)

var testConf = AlbumConfig{
	AlbumDir: "./tmp/",
}

func TestIndexUpdate(t *testing.T) {

	// Prepare test folder
	os.RemoveAll("./tmp/")
	utils.CreateFolder("./tmp/", 0775)
	defer os.RemoveAll("./tmp/")
	err := utils.RecursiveCopy("./testdata/", "./tmp/")
	if err != nil {
		log.Fatal(err)
	}

	// Index
	index, err := NewIndex(&testConf)
	if err != nil {
		panic(err)
	}

	index.UpdateAll()

	//pretty.Log(index.root)

	_, err = os.Open("./tmp/_albums.json")

	Convey("Album index file", t, func() {
		So(err, ShouldEqual, nil)
	})

	checkRoot(t, "Initial in memory index", &index.root)

	index.Shutdown()

	// Load back from disk and check it still matches
	index2, err := NewIndex(&testConf)
	if err != nil {
		panic(err)
	}

	checkRoot(t, "Initial in memory index", &index2.root)

	index2.Shutdown()
}

// check the initail root album
func checkRoot(t *testing.T, testTitle string, root *Album) {
	Convey(testTitle, t, func() {
		So(len(root.pics), ShouldEqual, 2)

		So(len(root.Children), ShouldEqual, 2)
		fb := root.Child("foobar")
		So(fb, ShouldEqual, nil)

		d1 := root.Child("d1")
		So(d1, ShouldNotEqual, nil)
		So(len(d1.Children), ShouldEqual, 1)
		So(len(d1.pics), ShouldEqual, 2)

		d2 := root.Child("d2")
		So(d2, ShouldNotEqual, nil)
		So(len(d2.Children), ShouldEqual, 0)
		So(len(d1.pics), ShouldEqual, 1)

		d11 := d1.Child("d11")
		So(d11, ShouldNotEqual, nil)
		So(len(d11.Children), ShouldEqual, 0)
		So(len(d1.pics), ShouldEqual, 2)

		So(d11.Path, ShouldEqual, "d11")
		So(d11.Name, ShouldEqual, "d11")
		So(d11.Ordering, ShouldEqual, 0)
		So(d11.Hidden, ShouldEqual, false)
		So(d11.dirty, ShouldEqual, true)

	})
}
