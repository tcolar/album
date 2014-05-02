// History: Apr 23 14 tcolar Creation

package album

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Index struct {
	conf *AlbumConfig
	root Album // In memory index
}

// NewIndex creates a new idexer
func NewIndex(conf *AlbumConfig) (indexer *Index, err error) {

	//filePath := path.Join(conf.AlbumDir, "_album.json")

	index := Index{
		conf: conf,
	}
	index.loadAlbumIndex()

	return &index, nil
}

func (i *Index) Shutdown() {
	// TODO: Anyhting ?
}

// Update scans the image folder and update the database with found album and pictures
// For new & updated pictures it will also create scaled down versions & thumbnails
func (i *Index) UpdateAll() {

	// Todo : load existing stuff from json

	dirtyAlbums := i.Cleanup(&i.root)
	dirtyAlbums = i.UpdateAlbum(i.conf.AlbumDir, &i.root) || dirtyAlbums

	if dirtyAlbums {
		i.saveAlbumIndex()
	}

	log.Print("Index update completed.")
}

// Save the album index to file (JSON enncoded)
// TODO : might want to use a channel / sync this.
func (i *Index) saveAlbumIndex() {
	file := path.Join(i.conf.AlbumDir, "_albums.json")
	log.Printf("Saving albums index to %s", file)
	b, err := json.Marshal(i.root)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(file, b, 0644)
}

func (i *Index) loadAlbumIndex() {
	file := path.Join(i.conf.AlbumDir, "_albums.json")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		i.root = Album{}
		return
	}
	log.Printf("Loading album index from %s", file)
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &i.root)
	if err != nil {
		panic(err)
	}
}

// UpdateAlbum removes albums & pics that are no longer present on disk, from the index.
func (i *Index) Cleanup(album *Album) bool {
	dirty := false
	//for _, child := range album.Children {
	// browse albums -> path gone -> remove from db (along wth pics)
	// if any album gone -> remove from index, dirty = true
	//}
	// TODO
	// browse images -> path gone -> remove from db + remove from highlights ?
	// delete scaled down image and thumbails for gone images ?
	// rewrite that album image list (json)
	return dirty
}

// UpdateAlbum rescursively scans the album file system and update the index.
// Returns wether the album structure chnaged (new/removed albums)
func (i *Index) UpdateAlbum(dir string, album *Album) bool {

	dirty := false

	/*rel, err := filepath.Rel(i.conf.AlbumDir, dir)

	  if err != nil {
	    panic(err)
	  }*/

	// Browse disk
	// -> new album (path) -> create
	// -> new image (path) -> create
	// -> Updated image (ts) -> ??
	// create scaled down image and thumbails for new images

	/* Upate pics of this album
	err = i.UpdatePics(dir, key)
	if err != nil {
	  return err
	}

	highlight := "" // TODO
	*/
	// If it's a new album then create it
	/*child := AlbumByPath(i, key)
	  if album == nil {
	    // Album represents a collection of Pics
	    parentPath := ""
	    if parent != nil {
	      parentPath = parent.Path
	    }
	    album = &Album{
	      Path:         key,
	      Parent:       parentPath,
	      Name:         rel,
	      Description:  "",
	      HighlightPic: highlight,
	      Ordering:     0,
	      Hidden:       false,
	    }
	    album.Save(i)
	  }*/

	// Recurse into subalbums
	files, _ := ioutil.ReadDir(dir)

	for _, f := range files {
		nm := f.Name()
		fp := path.Join(dir, nm)
		if f.IsDir() {
			child := album.Child(nm)
			if child == nil {
				dirty = true
				// create it
				child = &Album{
					Path:         nm,
					Name:         nm,
					Description:  "",
					HighlightPic: "", // TODO
					Ordering:     0,
					Hidden:       false,
					dirty:        true,
				}
				dirty = i.UpdateAlbum(fp, child) || dirty
				album.Children = append(album.Children, *child)
			} else {
				dirty = i.UpdateAlbum(fp, child) || dirty
			}
		}
	}
	return dirty
}

/*func (i *Index) UpdatePics(dir string, albumKey string) (err error) {
  rel, err := filepath.Rel(i.conf.AlbumDir, dir)

  if err != nil {
    return err
  }

  files, _ := ioutil.ReadDir(dir)
  for _, f := range files {
    if !f.IsDir() && IsImage(f) {
      key := i.normalizeKey("", path.Join(rel, f.Name()))
      pic := PicByPath(i, key)
      if pic == nil { // TODO: what if modtime changed ?
        pic = &Pic{
          Path:        key,
          Album:       albumKey,
          Name:        f.Name(),
          Description: "",
          ModTime:     f.ModTime().Unix(),
        }
        pic.Save(i) // -> No need to save as it has no "custom" data yet
      }
    }
  }
  return nil
}*/
