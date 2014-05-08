// History: Apr 23 14 tcolar Creation

package album

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
)

// TODO: Use channels to write json files

type Index struct {
	conf   *AlbumConfig
	root   Album // In memory tree index
	imgSvc ImageSvc
}

// NewIndex creates a new indexer
func NewIndex(conf *AlbumConfig) (indexer *Index, err error) {

	index := Index{
		conf:   conf,
		imgSvc: ImageSvc{},
	}

	index.loadAlbumIndex()

	index.loadAllAlbumPics(&index.root, index.conf.AlbumDir)

	return &index, nil
}

// UpdateAll scans the image folder and update the database with found album and pictures
// For new & updated pictures it will also create scaled down versions & thumbnails
func (i *Index) UpdateAll() {

	log.Print("Starting index update : %s", i.conf.AlbumDir)

	dirtyAlbums := i.Cleanup(&i.root)
	dirtyAlbums = i.UpdateAlbum(i.conf.AlbumDir, &i.root) || dirtyAlbums
	dirtyAlbums = i.UpdateHighLights(&i.root, "") || dirtyAlbums

	// This means the album tree structure itself has changed
	if dirtyAlbums {
		i.saveAlbumIndex()
	}

	// save any album whose content (pics) have changed
	i.saveDirtyAlbums(&i.root, i.conf.AlbumDir)

	log.Print("Index update completed.")
}

// Cleanup removes albums & pics that are no longer present on disk, from the index.
// Returns wether the album index is dirty
// Also will set dirty var on individual albums as needed
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

	// Recurse into subalbums
	files, _ := ioutil.ReadDir(dir)
	var mostRecentPicTs int64 = 0
	if p := album.LatestPic(); p != nil {
		mostRecentPicTs = p.ModTime
	}

	for _, f := range files {
		nm := f.Name()
		fp := path.Join(dir, nm)
		if f.IsDir() && nm != "_scaled" {
			// dir
			child := album.Child(nm)
			if child == nil {
				dirty = true
				child = &Album{
					Path:         nm,
					Name:         nm,
					HighlightPic: "",
				}
				dirty = i.UpdateAlbum(fp, child) || dirty
				album.Children = append(album.Children, *child)
			} else {
				dirty = i.UpdateAlbum(fp, child) || dirty
			}
		} else {
			// file
			ts := f.ModTime().Unix()
			if i.imgSvc.IsImage(f) && ts > mostRecentPicTs {
				album.dirty = true
				pic := Pic{
					Path:    f.Name(),
					Name:    f.Name(),
					ModTime: ts,
				}
				// TODO: start those in a single background go routine / queue (channel) ?
				err := i.createScaledImages(fp)
				if err != nil {
					log.Print(err)
				} else {
					album.pics = append(album.pics, pic)
				}
			}
		}
	}

	return dirty
}

//  createScaledImages creates scaled down version of the images (thumbnails etc..)
func (i *Index) createScaledImages(fp string) error {
	log.Printf("Creating scaled images for %s", fp)
	dest, err := i.scaledPath(fp, "thumb", ".png")
	if err != nil {
		return err
	}
	return i.imgSvc.CreateThumbnail(fp, dest, 200, 200)
}

// Returns the web path of a scaled image
func (i *Index) scaledPath(fp, prefix, ext string) (patht string, err error) {
	rel, err := filepath.Rel(i.conf.AlbumDir, fp)
	if err != nil {
		return "", err
	}
	file := filepath.Base(rel)
	file = file[:len(file)-len(path.Ext(file))] + ext

	dest := path.Join(i.conf.AlbumDir, "_scaled", prefix, filepath.Dir(rel), file)
	return dest, nil
}

// Recursively save all albums whose content is dirty
func (i *Index) saveDirtyAlbums(album *Album, dir string) {
	if album.dirty {
		log.Printf("sda %s %s", dir, album.Path)
		// Sort them before saving
		sort.Sort(album.pics)
		i.saveAlbumPics(album, dir)
	}
	// recurse
	for c, _ := range album.Children {
		child := &album.Children[c]
		i.saveDirtyAlbums(child, path.Join(dir, child.Name))
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
}

// Load the pictures of all albums from json (recursively)
func (i *Index) loadAllAlbumPics(album *Album, dir string) {
	i.loadAlbumPics(album, dir)
	for c, _ := range album.Children {
		subPtr := &album.Children[c] // need to do this to pass by reference
		i.loadAllAlbumPics(subPtr, path.Join(dir, subPtr.Name))
	}
}

// loadAlbumPics load the pictures of a given album from the json file
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
	if album.pics == nil {
		album.pics = Pics{}
	}
	err = json.Unmarshal(b, &album.pics)
	if err != nil {
		panic(err)
	}
}

// Returns wether the album index s drty
func (i *Index) UpdateHighLights(a *Album, dir string) bool {

	dirty := false

	// Recurse first since an album highlight might bubbe up.
	dir = path.Join(dir, a.Path)
	for c, _ := range a.Children {
		child := &a.Children[c]
		dirty = i.UpdateHighLights(child, dir) || dirty
	}

	if len(a.HighlightPic) > 0 && a.Pic(a.HighlightPic) != nil {
		// Valid highlight, leave it alone
		return dirty
	}

	// if no highlight defined return first pic of album
	if len(a.pics) > 0 {
		dirty = true
		p := a.pics[0].Path
		nm := p[:len(p)-len(filepath.Ext(p))] + ".png"
		a.HighlightPic = path.Join(dir, nm)
		return dirty
	}

	// If no images in album, return highlight of first subalbum
	if len(a.Children) > 0 {
		dirty = true
		a.HighlightPic = a.Children[0].HighlightPic
		return dirty
	}

	// Nothing found, leave it alone, will try again next time
	a.HighlightPic = ""

	return dirty

}
