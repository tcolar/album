// History: May 08 14 tcolar Creation

package album

// Storer, stores and retrieves the album & pics metadata somewhere (say a DB or such)
type Storer interface {
	// GetRoot returns the root album (nil if none)
	GetRoot() (*Album, error)

	GetAlbum(id string) (*Album, error)
	CreateAlbum(album *Album) error
	UpdateAlbum(album *Album) error

	// DelAbum deletes an album and all its subalbums ad pictures (recursively)
	DelAlbum(id string) error

	// GetAlbum return the children albums of this album
	GetAlbums(id string) ([]Album, error)

	// GetAlbumPics return the pictures of this album
	GetAlbumPics(id string) ([]Pic, error)

	GetPic(id string) (*Pic, error)
	CreatePic(pic *Pic) error
	UpdatePic(pic *Pic) error

	// DelPic remves a picture
	DelPic(id string) error

	// To be called when exiting to close all resources
	Shutdown() error
}
