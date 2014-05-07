// History: Apr 21 14 tcolar Creation

package main

import (
	"flag"
	"log"
	"os"

	"github.com/tcolar/album"
)

func main() {

	port := flag.Int("port", 3000, "Web service port")

	if len(os.Args) < 2 {
		log.Fatal("Expected album directory as first parameter.")
	}

	flag.Parse()

	conf := album.AlbumConfig{
		AlbumDir: os.Args[1],
		Port:     *port,
	}
	s := album.NewServer(conf)
	s.Run()
}
