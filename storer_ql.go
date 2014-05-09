// History: May 08 14 tcolar Creation

package album

import (
	"os"

	"log"

	"github.com/cznic/kv"
	"labix.org/v2/mgo/bson"
)

// KvStorer : Implementation of storer interface using the KV key/value store.
// Structures are serialized to bson.
type KvStorer struct {
	db *kv.DB
}

func NewKvStorer(dbPath string) (*KvStorer, error) {
	var db *kv.DB
	var err error
	if _, err = os.Stat(dbPath); err == nil {
		db, err = kv.Open(dbPath, &kv.Options{})
	} else {
		db, err = kv.Create(dbPath, &kv.Options{})
	}
	if err != nil {
		return nil, err
	}
	return &KvStorer{
		db: db,
	}, nil
}

func (s *KvStorer) GetRoot() (album *Album, e error) {
	key := []byte(rootName)
	root, err := s.GetAlbum(rootName)
	if err != nil {
		return nil, err
	}
	if root == nil {
		root = &Album{Id: rootName}
		bytes, err := s.encAlbum(root)
		if err != nil {
			return nil, err
		}
		s.db.Set(key, bytes)
		return root, s.CreateAlbum(root)
	}
	return root, nil
}

func (s *KvStorer) GetAlbum(id string) (album *Album, e error) {
	bytes, err := s.db.Get([]byte(id), nil)
	if err != nil {
		return nil, err
	}
	return s.decAlbum(bytes)
}
func (s *KvStorer) CreateAlbum(album *Album) error {
	return s.UpdateAlbum(album)
}

func (s *KvStorer) UpdateAlbum(album *Album) error {
	bytes, err := s.encAlbum(album)
	if err != nil {
		return err
	}
	log.Printf("%s -> %v", album.Id, bytes)
	return s.db.Set([]byte(album.Id), bytes)
}

func (s *KvStorer) DelAlbum(id string) error {
	return nil // TODO
}

func (s *KvStorer) GetAlbums(id string) (albums []Album, err error) {
	return []Album{}, nil // TODO
}

func (s *KvStorer) GetAlbumPics(id string) (pics []Pic, e error) {
	return []Pic{}, nil // TODO
}

func (s *KvStorer) GetPic(id string) (pic *Pic, e error) {
	return &Pic{}, nil // TODO
}
func (s *KvStorer) CreatePic(pic *Pic) error {
	return nil // TODO
}
func (s *KvStorer) UpdatePic(pic *Pic) error {
	return nil // TODO
}
func (s *KvStorer) DelPic(id string) error {
	return nil // TODO
}

func (s *KvStorer) Shutdown() error {
	log.Print("Closing DB")
	return s.db.Close()
}

// encAlbum encodes an album to bson []byte
func (s *KvStorer) encAlbum(album *Album) ([]byte, error) {
	return bson.Marshal(album)
}

// devAlbum decodes an album from bson []byte
func (s *KvStorer) decAlbum(bytes []byte) (*Album, error) {
	if bytes == nil {
		return nil, nil
	}
	album := &Album{}
	err := bson.Unmarshal(bytes, album)
	return album, err
}

// encPic encodes a album to bson []byte
func (s *KvStorer) encPic(pic *Pic) ([]byte, error) {
	return bson.Marshal(pic)
}

// devPic decodes a pic from bson []byte
func (s *KvStorer) decPic(bytes []byte) (*Pic, error) {
	if bytes == nil {
		return nil, nil
	}
	pic := &Pic{}
	err := bson.Unmarshal(bytes, pic)
	return pic, err
}

var rootName string = "/"
