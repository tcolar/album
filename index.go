// History: Apr 23 14 tcolar Creation

package album

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
)

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

// UpdateAll scans the image folder and update the database with found album and pictures
// For new & updated pictures it will also create scaled down versions & thumbnails
func (i *Index) UpdateAll() {

	log.Printf("Starting index update : %s", i.conf.AlbumDir)

	i.Cleanup(i.rootAlbum())
	i.UpdateAlbum(i.conf.AlbumDir, i.rootAlbum())
	i.UpdateHighLights(i.rootAlbum())

	log.Print("Index update completed.")
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
					Path:         id,
					Name:         nm,
					HighlightPic: "",
				}
			}
			log.Printf("Indexing Album %s", child.Id)
			err = i.store.CreateAlbum(child)
			if err != nil {
				log.Print(err)
				continue
			}
			i.UpdateAlbum(fp, child)
		} else {
			// file
			ts := f.ModTime().Unix()
			if i.imgSvc.IsImage(f) {
				id := path.Join(album.Id, nm)
				img, err := i.store.GetPic(id)
				if err != nil {
					log.Print(err)
					continue
				}
				if img != nil && img.ModTime >= ts {
					continue // alreday up to date
				}
				w, h, err := i.imgSvc.ImageSize(fp)
				if err != nil {
					log.Print(err)
					continue
				}
				pic := &Pic{
					Id:      id,
					AlbumId: album.Id,
					Path:    id,
					Name:    f.Name(),
					ModTime: ts,
					Width:   w,
					Height:  h,
				}
				err = i.createScaledImages(fp, w, h)
				if err != nil {
					log.Print(err)
				} else {
					log.Printf("Indexing Pic %s", pic.Id)
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

// UpdateHighLights recursively make sure all albums have an highlight
func (i *Index) UpdateHighLights(a *Album) {

	// Recurse first since an album highlight might bubbe up.
	subs := i.subAlbums(a)
	for c, _ := range subs {
		sub := &subs[c]
		i.UpdateHighLights(sub)
	}

	if len(a.HighlightPic) > 0 && i.pic(a.HighlightPic) != nil {
		// Valid highlight, leave it alone
		return
	}

	// if no highlight defined return first pic of album
	pics := i.albumPics(a)
	if len(pics) > 0 {
		p := pics[0].Path
		nm := p[:len(p)-len(filepath.Ext(p))] + ".png"
		a.HighlightPic = nm
		i.store.UpdateAlbum(a)
		return
	}

	// If no images in album, return highlight of first subalbum that has one
	for c, _ := range subs {
		sub := &subs[c]
		if len(sub.HighlightPic) > 0 {
			a.HighlightPic = sub.HighlightPic
			i.store.UpdateAlbum(a)
			return
		}
	}
	// Nothing found, leave it alone, will try again next time
}

func (i *Index) album(id string) *Album {
	album, err := i.store.GetAlbum(id)
	if err != nil {
		log.Fatal(err) // TODO: log & ignore ?
	}
	return album
}

func (i *Index) albumPics(album *Album) []Pic {
	pics, err := i.store.GetAlbumPics(album.Id)
	if err != nil {
		log.Fatal(err) // TODO: log & ignore ?
	}
	return pics
}

// createScaledImages creates scaled down version of the images (thumbnails etc..)
// Thumbnail, small, medium, large
// TODO: If the scaled version would be higher reolution than the original maybe
// just create a link (or copy on windows) rather than waste disk space ?
func (i *Index) createScaledImages(fp string, w, h int) error {
	log.Printf("Creating scaled images for %s", fp)
	c := i.conf

	for key, size := range c.MediaSizes {
		if w > size.MaxScaleWidth || h > size.MaxScaleHeight {
			ext := path.Ext(fp)
			if size.PadOnScale {
				ext = ".png"
			}
			dest, err := i.scaledPath(fp, key, ext)
			if err != nil {
				return err
			}

			err = i.imgSvc.CreateScaled(fp, dest, size.MaxScaleWidth, size.MaxScaleHeight, size.PadOnScale)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (i *Index) pic(id string) *Pic {
	pic, err := i.store.GetPic(id)
	if err != nil {
		log.Fatal(err) // TODO: log & ignore ?
	}
	return pic
}

func (i *Index) rootAlbum() *Album {
	root, err := i.store.GetRoot()
	if err != nil {
		log.Fatalf("Failed to get root album. %v", err)
	}
	return root
}

// scaledPath returns the web path of a scaled image
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

func (i *Index) subAlbums(album *Album) []Album {
	albums, err := i.store.GetAlbums(album.Id)
	if err != nil {
		log.Fatal(err) // TODO: log & ignore ?
	}
	return albums
}
