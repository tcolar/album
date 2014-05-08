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

func NewServer(conf AlbumConfig) *Server {

	CreateAdminUser(conf.AdminPassword)

	return &Server{
		Conf: conf,
	}
}

type Server struct {
	Conf  AlbumConfig
	index *Index
	api   Api
}

func (s *Server) Run() {

	index, err := NewIndex(&s.Conf)
	if err != nil {
		panic(err)
	}
	s.index = index
	s.api = Api{}

	// Index all albums & images asynchronously
	go index.UpdateAll()

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
	http.ListenAndServe(fmt.Sprintf(":%d", s.Conf.Port), m)
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
	log.Print(parts)
	album := &s.index.root
	log.Print(albums)
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
