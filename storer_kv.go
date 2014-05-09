// History: May 08 14 tcolar Creation

package album

import (
  "fmt"
  "os"
  "strings"

  "log"

  "github.com/cznic/kv"
  "labix.org/v2/mgo/bson"
)

// KvStorer : Implementation of album.Storer interface backed by KV key/value store.
// Values are structures serialized to bson.
type KvStorer struct {
  db *kv.DB
}

// NewKvStorer constructor for a KV stored back Storer implememtation
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

func (s *KvStorer) GetRoot() (*Album, error) {
  key := []byte(s.albumKey(rootName))
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

func (s *KvStorer) GetAlbum(id string) (*Album, error) {
  // Get is (value, key) while Set is (key, value) >:<
  bytes, err := s.db.Get(nil, []byte(s.albumKey(id)))
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
  return s.db.Set([]byte(s.albumKey(album.Id)), bytes)
}

func (s *KvStorer) DelAlbum(id string) error {
  // recurse in subalbums
  subAlbums, err := s.GetAlbums(id)
  if err != nil {
    return err
  }
  for _, sub := range subAlbums {
    s.DelAlbum(sub.Id)
  }
  // delete album pics
  pics, err := s.GetAlbumPics(id)
  if err != nil {
    return err
  }
  for _, pic := range pics {
    s.DelPic(pic.Id)
  }
  // remove album itself
  return s.db.Delete([]byte(s.albumKey(id)))
}

func (s *KvStorer) GetAlbums(id string) ([]Album, error) {
  enum, _, err := s.db.Seek([]byte(s.albumKey(id)))
  if err != nil {
    return nil, err
  }
  albums := []Album{}
  k, v, e := enum.Next()
  key := string(k)
  kid := s.albumKey(id)
  for ; e == nil && strings.HasPrefix(key, kid); k, v, e = enum.Next() {
    key = string(k)
    if key == id {
      continue // First key will be the album itself, skipping it
    }
    album, err := s.decAlbum(v)
    if err == nil && album.ParentId == id {
      albums = append(albums, *album)
    }
  }
  return albums, nil
}

func (s *KvStorer) GetAlbumPics(id string) ([]Pic, error) {
  enum, _, err := s.db.Seek([]byte(s.picKey(id)))
  if err != nil {
    return nil, err
  }
  pics := []Pic{}
  k, v, e := enum.Next()
  key := string(k)
  kid := s.picKey(id)
  for ; e == nil && strings.HasPrefix(key, kid); k, v, e = enum.Next() {
    key = string(k)
    pic, err := s.decPic(v)
    if err == nil && pic.AlbumId == id {
      pics = append(pics, *pic)
    }
  }
  return pics, nil
}

func (s *KvStorer) GetPic(id string) (*Pic, error) {
  // Get is (value, key) while Set is (key, value) >:<
  bytes, err := s.db.Get(nil, []byte(s.picKey(id)))
  if err != nil {
    return nil, err
  }
  return s.decPic(bytes)
}

func (s *KvStorer) CreatePic(pic *Pic) error {
  return s.UpdatePic(pic)
}

func (s *KvStorer) UpdatePic(pic *Pic) error {
  bytes, err := s.encPic(pic)
  if err != nil {
    return err
  }
  return s.db.Set([]byte(s.picKey(pic.Id)), bytes)
}

func (s *KvStorer) DelPic(id string) error {
  return s.db.Delete([]byte(s.picKey(id)))
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

// albumKey  :Internal KV store album key ("A:key")
func (s *KvStorer) albumKey(key string) string {
  return fmt.Sprintf("A:%s", key)
}

// albumKey  :Internal KV store pic key ("I:key")
func (s *KvStorer) picKey(key string) string {
  return fmt.Sprintf("I:%s", key)
}

var rootName string = "/"