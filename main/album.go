// History: Apr 21 14 tcolar Creation

package main

import (
	"log"
	"os"

	"github.com/tcolar/album"
)

func main() {

	//api := flag.Bool("api", true, "Enable API")

	if len(os.Args) < 2 {
		log.Fatal("Expected album directory as first parameter.")
	}

	conf := album.AlbumConfig{
		AlbumDir: os.Args[1],
	}
	s := album.NewServer(conf)
	s.Run()
}
