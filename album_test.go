// History: Apr 26 14 tcolar Creation

package album

import (
	"log"
	"os"
	"testing"

	"github.com/cznic/kv"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tcolar/utils"
)

var testConf = AlbumConfig{
	AlbumDir: "./tmp/",
}

func TestIndex(t *testing.T) {

	// Prepare test folder
	os.RemoveAll("./tmp/")
	utils.CreateFolder("./tmp/", 0775)
	defer os.RemoveAll("./tmp/")
	err := utils.RecursiveCopy("./testdata/", "./tmp/")
	if err != nil {
		log.Fatalf("Failed to get recursive copy. %v", err)
	}

	// Index
	store, err := NewKvStorer("./tmp/album_db")
	if err != nil {
		log.Fatalf("Failed to open db. %v", err)
	}
	index, err := NewIndex(&testConf, store)
	if err != nil {
		panic(err)
	}

	sz, err := store.db.Size()
	log.Printf("Size: %d %v", sz, err)

	index.UpdateAll()

	sz, err = store.db.Size()
	log.Printf("Size: %d %v", sz, err)
	b, err := store.db.Get([]byte("/d1"), nil)
	log.Print("Got: %s", b)
	dumpDb(kv.DB(*store.db))

	// TODO: check no changes
	checkRoot(t, "Initial index", index)

	// TODO: Test cleanup, updates, chnages etc ....
	// TODO: Test ordering
}

// check the initial root album
func checkRoot(t *testing.T, testTitle string, index *Index) {

	Convey(testTitle, t, func() {
		root, err := index.store.GetRoot()
		if err != nil {
			log.Fatal(err)
		}

		pics := index.AlbumPics(root)
		albums := index.AlbumPics(root)

		So(root, ShouldNotEqual, nil)

		So(len(pics), ShouldEqual, 2)
		So(len(albums), ShouldEqual, 2)

		/*fb := root.Child("foobar")
		  So(fb, ShouldEqual, nil)

		  So(len(root.pics), ShouldEqual, 2)
		  So(root.Pic("Cart_1.png"), ShouldEqual, nil)
		  So(root.Pic("Minus.png"), ShouldNotEqual, nil)

		  d1 := root.Child("d1")

		  So(d1, ShouldNotEqual, nil)
		  So(len(d1.Children), ShouldEqual, 1)
		  So(len(d1.pics), ShouldEqual, 0)

		  d2 := root.Child("d2")
		  So(d2, ShouldNotEqual, nil)
		  So(len(d2.Children), ShouldEqual, 0)
		  So(len(d2.pics), ShouldEqual, 1)

		  d11 := d1.Child("d11")
		  So(d11, ShouldNotEqual, nil)
		  So(len(d11.Children), ShouldEqual, 0)
		  So(len(d11.pics), ShouldEqual, 2)

		  So(d11.Path, ShouldEqual, "d11")
		  So(d11.Name, ShouldEqual, "d11")
		  So(d11.Ordering, ShouldEqual, 0)
		  So(d11.Hidden, ShouldEqual, false)*/
	})
}

func dumpDb(db kv.DB) {
	enum, err := db.SeekFirst()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Db keys:")
	i := 0
	for k, _, e := enum.Next(); e == nil && i < 10; i++ {
		log.Print(string(k))
	}
}
