// History: Apr 21 14 tcolar Creation

package main

import "github.com/tcolar/album"

func main() {
	conf := album.AlbumConfig{
		AlbumDir: "/home/tcolar/albums/",
	}
	s := album.Server{Conf: conf}
	s.Run()
}
