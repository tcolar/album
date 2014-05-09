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

	index.UpdateAll()

	//dumpDb(kv.DB(*store.db))

	// TODO: check no changes
	checkRoot(t, "Initial index", index)

	// TODO: Test cleanup, updates, chnages etc ....
	// TODO: Test ordering

	Convey("Checking highlights.", t, func() {
		So(album(index, "/d1/d11").HighlightPic, ShouldEqual, "/d1/d11/Cart_1.png")
		So(album(index, "/d1").HighlightPic, ShouldEqual, "/d1/d11/Cart_1.png")
		So(album(index, "/d2").HighlightPic, ShouldEqual, "/d2/DownArrow.png")
		So(album(index, "/").HighlightPic, ShouldEqual, "/Minus.png")
	})

	Convey("Deleting image from index should work.", t, func() {
		key := "/d1/d11/Cart_2.png"
		So(pic(index, key), ShouldNotEqual, nil)
		err = index.store.DelPic(key)
		if err != nil {
			panic(err)
		}
		So(pic(index, key), ShouldEqual, nil)
		So(pic(index, "/d1/d11/Cart_1.png"), ShouldNotEqual, nil) // Should still be there
	})

	Convey("Deleting album from index should work.", t, func() {
		key := "/d1"
		So(album(index, key), ShouldNotEqual, nil)
		err = index.store.DelAlbum(key)
		if err != nil {
			panic(err)
		}
		So(album(index, key), ShouldEqual, nil)
		So(album(index, "/d1/d11"), ShouldEqual, nil)          // subAlbums should be gone too
		So(pic(index, "/d1/d11/Cart_1.png"), ShouldEqual, nil) // and any pics within
		So(album(index, "/d2"), ShouldNotEqual, nil)           // should still be there
	})
}

// check the initial root album
func checkRoot(t *testing.T, testTitle string, index *Index) {

	Convey(testTitle, t, func() {
		root, err := index.store.GetRoot()
		if err != nil {
			log.Fatal(err)
		}

		// checking root album
		So(root, ShouldNotEqual, nil)
		So(album(index, "/"), ShouldNotEqual, nil)
		So(album(index, "/").Id, ShouldEqual, "/")

		pics := index.albumPics(root)
		albums := index.subAlbums(root)

		So(len(pics), ShouldEqual, 2)
		So(len(albums), ShouldEqual, 2)

		So(album(index, "/foobar"), ShouldEqual, nil)

		So(pics[0].Id, ShouldEqual, "/Minus.png")
		So(pics[1].Id, ShouldEqual, "/Plus.png")

		d1 := album(index, "/d1")

		So(d1, ShouldNotEqual, nil)
		So(len(index.subAlbums(d1)), ShouldEqual, 1)
		So(len(index.albumPics(d1)), ShouldEqual, 0)

		d2 := album(index, "/d2")
		So(d2, ShouldNotEqual, nil)
		So(len(index.subAlbums(d2)), ShouldEqual, 0)
		So(len(index.albumPics(d2)), ShouldEqual, 1)

		d11 := album(index, "/d1/d11")
		So(d11, ShouldNotEqual, nil)
		So(len(index.subAlbums(d11)), ShouldEqual, 0)
		So(len(index.albumPics(d11)), ShouldEqual, 2)

		So(d11.Id, ShouldEqual, "/d1/d11")
		So(d11.ParentId, ShouldEqual, "/d1")
		So(d11.Path, ShouldEqual, "/d1/d11")
		So(d11.Name, ShouldEqual, "d11")
		So(d11.Ordering, ShouldEqual, 0)
		So(d11.Hidden, ShouldEqual, false)

		So(pic(index, "/Cart_1.png"), ShouldEqual, nil)
		cart1 := pic(index, "/d1/d11/Cart_1.png")
		So(cart1, ShouldNotEqual, nil)
		So(cart1.Id, ShouldEqual, "/d1/d11/Cart_1.png")
		So(cart1.AlbumId, ShouldEqual, "/d1/d11")
		So(cart1.Path, ShouldEqual, "/d1/d11/Cart_1.png")
		So(cart1.Name, ShouldEqual, "Cart_1.png")
		So(cart1.Ordering, ShouldEqual, 0)
		So(cart1.Hidden, ShouldEqual, false)
	})
}

func album(index *Index, key string) *Album {
	album, err := index.store.GetAlbum(key)
	if err != nil {
		log.Fatal(err)
	}
	return album
}

func pic(index *Index, key string) *Pic {
	pic, err := index.store.GetPic(key)
	if err != nil {
		log.Fatal(err)
	}
	return pic
}

func dumpDb(db kv.DB) {
	enum, err := db.SeekFirst()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Db keys:")
	k, _, e := enum.Next()
	for ; e == nil; k, _, e = enum.Next() {
		log.Print(string(k))
	}
}
