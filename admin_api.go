// History: May 06 14 tcolar Creation

package album

import (
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/sessionauth"
)

// API provides some API's to the image server
// TODO: Albums, Pics API's
type Api struct{}

var ApiVersion string = "1.0.0"

// initRoutes add this API's routes to martini
func (a Api) initRoutes(m *martini.ClassicMartini) {

	m.Get("/_api/v1/img/version", a.Version)

	m.Post("/_api/v1/img/rotate", sessionauth.LoginRequired, binding.Bind(JsonData{}), a.ImageRotate)

	m.Post("/_api/v1/img/delete", sessionauth.LoginRequired, binding.Bind(JsonData{}), a.ImageDelete)
}

// /_api/v1/img/version returns the current API version
func (a Api) Version() string {
	return ApiVersion
}

// /_api/v1/img/rotate rotates & replace the given picture (by path)
// Also rotates any scaled versions of it (thubnails etc...)
// Json Params:
//   - ImagePath : full path to the image (from album root)
//   - Angle : Rotation angle in degrees (ie: 90)
func (a Api) ImageRotate(res http.ResponseWriter, data JsonData) {
	//s.servePics(r, req, res)
	log.Print("Rotate TBD")
	// todo : rotate and save
}

// /_api/v1/img/delete completely removes an image (for good).
// Also removes any scaled versions of it (thubnails etc...)
// Json Params:
//   - ImagePath : full path to the image (from album root)
func (a Api) ImageDelete(res http.ResponseWriter, data JsonData) {
	//s.servePics(r, req, res)
	log.Print("Delete TBD")
	// todo : delete file
	// remove from album index (pics) -> in memory / on file too ?
}

type JsonData struct {
	ImagePath string
	Angle     int
}
