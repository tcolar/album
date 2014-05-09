// History: Apr 23 14 tcolar Creation

package album

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
)

// TODO: Use channels to write json files

type Index struct {
	conf   *AlbumConfig
	imgSvc ImageSvc
	store  Storer
}

// NewIndex creates a new indexer
func NewIndex(conf *AlbumConfig, store Storer) (indexer *Index, err error) {

	index := Index{
		conf:   conf,
		store:  store,
		imgSvc: ImageSvc{},
	}

	return &index, nil
}

func (i *Index) rootAlbum() *Album {
	root, err := i.store.GetRoot()
	if err != nil {
		log.Fatalf("Failed to get root album. %v", err)
	}
	return root
}

// UpdateAll scans the image folder and update the database with found album and pictures
// For new & updated pictures it will also create scaled down versions & thumbnails
func (i *Index) UpdateAll() {

	log.Printf("Starting index update : %s", i.conf.AlbumDir)

	i.Cleanup(i.rootAlbum())
	i.UpdateAlbum(i.conf.AlbumDir, i.rootAlbum())
	//i.UpdateHighLights(i.rootAlbum(), "")

	log.Print("Index update completed.")
}

// Cleanup removes albums & pics that are no longer present on disk, from the index.
// Returns whether the album index is dirty
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
func (i *Index) UpdateAlbum(dir string, album *Album) {

	// Recurse into subalbums
	files, _ := ioutil.ReadDir(dir)

	for _, f := range files {
		nm := f.Name()
		fp := path.Join(dir, nm)
		if f.IsDir() && nm != "_scaled" {
			// dir
			id := path.Join(album.Id, nm)
			child, err := i.store.GetAlbum(id)
			if err != nil {
				log.Fatalf("Failed to get subalbum. %v", err)
			}
			if child == nil {
				child = &Album{
					Id:           id,
					ParentId:     album.Id,
					Path:         nm,
					Name:         nm,
					HighlightPic: "",
				}
			}
			log.Printf("Indexing Album %v", child)
			err = i.store.CreateAlbum(child)
			if err != nil {
				log.Print(err)
				continue
			}
			i.UpdateAlbum(fp, child)
		} else {
			// file
			ts := f.ModTime().Unix()
			if i.imgSvc.IsImage(f) && ts > album.LastPicModTime {
				id := path.Join(album.Id, nm)
				pic := &Pic{
					Id:      id,
					AlbumId: album.Id,
					Path:    f.Name(),
					Name:    f.Name(),
					ModTime: ts,
				}
				err := i.createScaledImages(fp)
				if err != nil {
					log.Print(err)
				} else {
					log.Printf("Indexing Pic %v", pic)
					err = i.store.CreatePic(pic)
					if err != nil {
						log.Print(err)
						continue
					}
				}
			}
		}
	}
}

func (i *Index) SubAlbums(album *Album) []Album {
	albums, err := i.store.GetAlbums(album.Id)
	if err != nil {
		log.Fatal(err) // TODO: log & ignore ?
	}
	return albums
}

func (i *Index) Album(id string) *Album {
	album, err := i.store.GetAlbum(id)
	if err != nil {
		log.Fatal(err) // TODO: log & ignore ?
	}
	return album
}

func (i *Index) AlbumPics(album *Album) []Pic {
	pics, err := i.store.GetAlbumPics(album.Id)
	if err != nil {
		log.Fatal(err) // TODO: log & ignore ?
	}
	return pics
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

// Returns whether the album index is drty
/*func (i *Index) UpdateHighLights(a *Album, dir string) bool {

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
}*/
