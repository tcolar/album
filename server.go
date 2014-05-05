// History: Apr 21 14 tcolar Creation

package album

import (
	"net/http"
	"path"
	"strings"

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

	// Index all albums & images asynchronously
	go index.UpdateAll()

	m := martini.Classic()
	m.Use(render.Renderer())

	root := conf.AlbumDir

	// To serve pictures
	// TODO: Protect json files ?
	m.Use(martini.Static(root, martini.StaticOptions{}))

	m.Get("/**", func(r render.Render, req *http.Request, res http.ResponseWriter) {

		albums := []Album{}
		pics := [][]string{}

		parts := strings.Split(req.URL.Path, "/")
		album := &index.root
		for _, p := range parts {
			if len(p) == 0 {
				continue
			}
			album = album.Child(p)
			if album == nil {
				break
			}
		}

		if album != nil {
			for _, a := range album.Children {
				albums = append(albums, a)
			}
			for _, p := range album.pics {
				pics = append(pics, []string{
					path.Join(req.URL.Path, p.Name),
					path.Join("/_scaled", "thumb", req.URL.Path, p.Name)})
			}
		}
		data := map[string]interface{}{
			"albums": albums,
			"pics":   pics,
		}
		r.HTML(200, "home", data)
	})

	m.Run()
}
