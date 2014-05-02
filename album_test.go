// History: Apr 26 14 tcolar Creation

package album

import "testing"

var testConf = AlbumConfig{
	AlbumDir: "/home/tcolar/albums/",
}

func TestIndex(t *testing.T) {
	index := NewIndex(conf)

	index.Update()
}