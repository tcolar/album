// History: Apr 21 14 tcolar Creation

package album

import (
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

// temp test config mock-up
var conf = AlbumConfig{
	AlbumDir: "/home/tcolar/albums/",
}

func Run() {

	index, err := NewIndex(&conf)
	if err != nil {
		panic(err)
	}

	// Index all abum & images asynchronously
	go index.UpdateAll()

	m := martini.Classic()
	m.Use(render.Renderer())

	root := conf.AlbumDir
	m.Use(martini.Static(root, martini.StaticOptions{
		Prefix: "_img",
	}))

	m.Get("/**", func(r render.Render, req *http.Request) {
		albums := []string{}
		pics := []string{}

		url := req.URL.Path[1:]
		if len(url) > 0 && url[len(url)-1] == '/' {
			url = url[0 : len(url)-1]
		}
		log.Printf("url: %s", url)
		/*album := index.Albums[url]
		  log.Printf("album: %v", album)
		  for _, pic := range album.Pics {
		    pics = append(pics, pic.Path)
		  }
		  for k, a := range index.Albums {
		    log.Printf("%s -> %s | %s", k, album.Path, a.ParentPath)
		    if album.Path == a.ParentPath {
		      albums = append(albums, a.Name)
		    }
		  }*/
		data := map[string]interface{}{
			"albums": albums,
			"pics":   pics,
		}
		r.HTML(200, "home", data)
	})

	defer func() {
		index.Shutdown()
	}()

	m.Run()
}
