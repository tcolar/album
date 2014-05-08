// History: May 07 14 tcolar Creation

package album

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/martini-contrib/sessionauth"
)

var users map[int64]User = map[int64]User{}

// TODO: This is temporary, need to create a password has ans such.
type User struct {
	Id            int64
	Username      string `form:"name"`
	Password      string `form:"password"`
	authenticated bool
}

// Login will preform any actions that are required to make a user model
// officially authenticated.
func (u *User) Login() {
	u.authenticated = true
}

// Logout will preform any actions that are required to completely
// logout a user.
func (u *User) Logout() {
	u.authenticated = false
}

func (u *User) IsAuthenticated() bool {
	return u.authenticated
}

func (u *User) UniqueId() interface{} {
	return u.Id
}

// GetById will populate a user object from a database model with
// a matching id.
func (u *User) GetById(id interface{}) error {
	var u2 User
	u2, ok := users[id.(int64)]
	if !ok {
		return fmt.Errorf("No user found for id : %v.", id)
	}
	u.Id = u2.Id
	u.Username = u2.Username
	u.Password = u2.Password
	return nil
}

// NewAdminUser creates the admin user with the given plaintext password
// TODO: secure this, use a password hash
func CreateAdminUser(password string) sessionauth.User {
	u := User{
		Id:       1,
		Username: "admin",
		Password: password,
	}
	if len(u.Password) == 0 {
		u.randPassword()
		log.Printf("*****Created a random admin password: '%s' ******", u.Password)
	}
	users[1] = u
	return &u
}

// Sets a random password
func (u *User) randPassword() {
	rand.Seed(time.Now().Unix())
	rnd := rand.Intn(990000) + 10000
	u.Password = strconv.Itoa(rnd)
}

func NewGuestUser() sessionauth.User {
	u := User{
		Id:       0,
		Username: "guest",
		Password: "",
	}
	users[0] = u
	return &u
}
