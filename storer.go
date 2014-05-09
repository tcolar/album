// History: May 08 14 tcolar Creation

package album

// Storer, stores and retrieves the album & pics metadata somewhere (say a DB or such)
type Storer interface {
	// GetRoot returns the root album (nil if none)
	GetRoot() (album *Album, err error)

	GetAlbum(id string) (album *Album, err error)
	CreateAlbum(album *Album) error
	UpdateAlbum(album *Album) error
	DelAlbum(id string) error

	// GetAlbum return the children albums of this album
	GetAlbums(id string) (albums []Album, err error)
	// GetAlbumPics return the pictures of this album
	GetAlbumPics(id string) (pics []Pic, err error)

	GetPic(id string) (pic *Pic, e error)
	CreatePic(pic *Pic) error
	UpdatePic(pic *Pic) error
	DelPic(id string) error

	// To be called when exiting to close all resources
	Shutdown() error
}
