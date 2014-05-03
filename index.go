// History: Apr 23 14 tcolar Creation

package album

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
)

type Index struct {
	conf *AlbumConfig
	root Album // In memory index
}

// NewIndex creates a new indexer
func NewIndex(conf *AlbumConfig) (indexer *Index, err error) {

	index := Index{
		conf: conf,
	}
	index.loadAlbumIndex()

	return &index, nil
}

// Update scans the image folder and update the database with found album and pictures
// For new & updated pictures it will also create scaled down versions & thumbnails
func (i *Index) UpdateAll() {

	// Todo : load existing stuff from json

	dirtyAlbums := i.Cleanup(&i.root)
	dirtyAlbums = i.UpdateAlbum(i.conf.AlbumDir, &i.root) || dirtyAlbums

	// This means the album tree structure itself has chnaged
	if dirtyAlbums {
		i.saveAlbumIndex()
	}

	// save any album whose content (pics) have changed
	i.saveDirtyAlbums(&i.root, i.conf.AlbumDir)

	log.Print("Index update completed.")
}

// UpdateAlbum removes albums & pics that are no longer present on disk, from the index.
// Returns wether the album index is dirty
// Also sets dirty bits on individual albums as needed
func (i *Index) Cleanup(album *Album) bool {
	dirty := false
	//for _, child := range album.Children {
	// browse albums -> path gone -> remove from album index
	// set parent to dirty
	//}
	// TODO
	// browse images -> path gone -> remove from db + remove from highlights ?
	// delete scaled down image and thumbails for gone images ?
	// rewrite that album image list (json)

	// TODO mae sure to mark album dirty if any images chnaged.

	return dirty
}

// UpdateAlbum rescursively scans the album file system and update the index.
// Returns wether the album structure chnaged (new/removed albums)
func (i *Index) UpdateAlbum(dir string, album *Album) bool {

	dirty := false

	// Browse disk
	// -> new album (path) -> create
	// -> new image (path) -> create
	// -> Updated image (ts) -> ??
	// create scaled down image and thumbails for new images

	//highlight := "" // TODO

	// Recurse into subalbums
	files, _ := ioutil.ReadDir(dir)
	var mostRecentPicTs int64 = 0
	if p := album.LatestPic(); p != nil {
		mostRecentPicTs = p.ModTime
	}

	for _, f := range files {
		nm := f.Name()
		fp := path.Join(dir, nm)
		if f.IsDir() {
			// dir
			child := album.Child(nm)
			if child == nil {
				dirty = true
				child = &Album{
					Path:         nm,
					Name:         nm,
					HighlightPic: "", // TODO
				}
				dirty = i.UpdateAlbum(fp, child) || dirty
				album.Children = append(album.Children, *child)
			} else {
				dirty = i.UpdateAlbum(fp, child) || dirty
			}
		} else {
			// file
			ts := f.ModTime().Unix()
			if IsImage(f) && ts > mostRecentPicTs {
				album.dirty = true
				pic := &Pic{
					Path:    f.Name(),
					Name:    f.Name(),
					ModTime: ts,
				}
				album.pics = append(album.pics, *pic)
			}
		}
	}

	return dirty
}

// Recursively save all albums whose content is dirty
func (i *Index) saveDirtyAlbums(album *Album, dir string) {
	if album.dirty {
		// Sort them before saving
		sort.Sort(album.pics)
		i.saveAlbumPics(album, dir)
		// TODO: switch back to ! dirty
	}
	// recurse
	for _, a := range album.Children {
		i.saveDirtyAlbums(&a, path.Join(dir, a.Name))
	}
}

func (i *Index) saveAlbumPics(album *Album, dir string) {
	file := path.Join(dir, "_pics.json")
	log.Printf("Saving albums index of %s to %s", album.Name, file)
	b, err := json.Marshal(album.pics)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(file, b, 0644)
}

// saveAlbumIndex saves the album index to file (JSON enncoded)
// TODO : might want to use a channel / sync this.
func (i *Index) saveAlbumIndex() {
	file := path.Join(i.conf.AlbumDir, "_albums.json")
	log.Printf("Saving album pics of  %s", file)
	b, err := json.Marshal(i.root)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(file, b, 0644)
}

// loadAlbumIndex load back the album index from file/json into memory
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
	// TODO : load pics
}

// loadAlbumPics load the pictures of a gven index from the json file
func (i *Index) loadAlbumPics(album *Album, dir string) {
	file := path.Join(dir, "_pics.json")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		album.pics = Pics{}
		return
	}
	log.Printf("Loading pics index of %s from %s", album.Name, file)
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &album.pics)
	if err != nil {
		panic(err)
	}
}
