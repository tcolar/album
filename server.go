// History: Apr 21 14 tcolar Creation

package album

import (
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func NewServer(conf AlbumConfig) *Server {
	return &Server{
		Conf: conf,
	}
}

type Server struct {
	Conf  AlbumConfig
	index *Index
}

func (s *Server) Run() {

	index, err := NewIndex(&s.Conf)
	if err != nil {
		panic(err)
	}
	s.index = index

	// Index all albums & images asynchronously
	go index.UpdateAll()

	m := martini.Classic()
	m.Use(render.Renderer())

	root := s.Conf.AlbumDir

	// To serve pictures
	// TODO: Protect json files ?
	m.Use(martini.Static(root, martini.StaticOptions{}))

	m.Get("/**", func(r render.Render, req *http.Request, res http.ResponseWriter) {
		s.servePics(r, req, res)
	})

	m.Run()
}

func (s *Server) servePics(r render.Render, req *http.Request, res http.ResponseWriter) {

	albums := []Album{}
	pics := [][]string{}

	parts := strings.Split(req.URL.Path, "/")
	album := &s.index.root
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
			nm := p.Path[:len(p.Path)-len(filepath.Ext(p.Path))] + ".png"
			pics = append(pics, []string{
				path.Join(req.URL.Path, p.Path),
				path.Join("/_scaled", "thumb", req.URL.Path, nm),
			})
		}
	}
	data := map[string]interface{}{
		"albums": albums,
		"pics":   pics,
	}
	r.HTML(200, "home", data)
}
