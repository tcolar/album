// History: Apr 21 14 tcolar Creation

package album

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
)

func NewServer(conf *AlbumConfig, store Storer) *Server {

	CreateAdminUser(conf.AdminPassword)

	index, err := NewIndex(conf, store)
	if err != nil {
		panic(err)
	}
	return &Server{
		Conf:  conf,
		index: index,
		api:   Api{},
	}
}

type Server struct {
	Conf  *AlbumConfig
	index *Index
	api   Api
}

func (s *Server) Run() {

	defer s.index.store.Shutdown() // Todo: use sync.waitgroup instead ?

	// Index all albums & images asynchronously
	go s.index.UpdateAll()

	m := martini.Classic()

	m.Use(render.Renderer())
	// TODO: Secure / random cookie store
	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("my_session", store))
	m.Use(sessionauth.SessionUser(NewGuestUser))
	sessionauth.RedirectUrl = "/_login"
	sessionauth.RedirectParam = "post-login"

	root := s.Conf.AlbumDir

	// Serve pictures & static content
	// TODO: Change this maybe to only server actual images ?
	m.Use(martini.Static(root, martini.StaticOptions{}))

	// Login
	m.Get("/_login", func(r render.Render) {
		r.HTML(200, "login", nil)
	})
	m.Post("/_login", binding.Bind(User{}), s.login)

	// API's
	s.api.initRoutes(m)

	// Album & pics pages
	m.Get("/**", s.servePics)

	log.Printf("Started on port %d", s.Conf.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.Conf.Port), m)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Done")
}

func (s *Server) login(session sessions.Session, user User, r render.Render, req *http.Request) {
	// For now can only login as the admin user
	admin := User{}
	admin.GetById(int64(1))
	if user.Username == admin.Username && admin.Password == user.Password {
		err := sessionauth.AuthenticateSession(session, &admin)
		if err != nil {
			r.JSON(500, err)
		}
		log.Printf("Logged as %v", admin)
	}

	params := req.URL.Query()
	redirect := params.Get(sessionauth.RedirectParam)
	r.Redirect(redirect)
	return
}

func (s *Server) servePics(r render.Render, req *http.Request, res http.ResponseWriter) {

	albums := []Album{}
	pics := [][]string{}

	parts := strings.Split(req.URL.Path, "/")
	id := fmt.Sprintf("/%s", path.Join(parts...))
	log.Printf("Id: %s", id)
	album, err := s.index.store.GetAlbum(id)
	if err != nil {
		log.Fatalf("Failed to get album. %v", err)
	}

	if album != nil {
		for _, a := range s.index.subAlbums(album) {
			albums = append(albums, a)
		}
		for _, p := range s.index.albumPics(album) {
			nm := p.Path[:len(p.Path)-len(filepath.Ext(p.Path))] + ".png"
			pics = append(pics, []string{
				p.Path,
				path.Join("/_scaled", "thumb", nm),
			})
		}
	}
	data := map[string]interface{}{
		"albums": albums,
		"pics":   pics,
	}
	r.HTML(200, "home", data)
}
