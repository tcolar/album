// History: Apr 26 14 tcolar Creation

package album

// Album represents a collection of Pics
type Album struct {
	Path         string // Relative (file) path from parent
	Name         string // Pretty name, might or not be same as Path
	Description  string
	HighlightPic string // Album highligh picture relative path
	Ordering     int    // if 0 will order by path
	Hidden       bool   // wether that album/folder should be hidden
	Children     []Album

	// Not serialized along with album structure
	pics  []Pic
	dirty bool // Wheter it's "dirty" (pics changed)
}

// Pic representents an individual picture / file
type Pic struct {
	Path        string
	Name        string
	Description string
	Ordering    int  // if 0 will order by path
	Hidden      bool // wether that picture should be hidden / desactivated
	ModTime     int64
}

// Album return a child album album by path, or nil if none found.
func (a Album) Child(path string) *Album {
	for _, child := range a.Children {
		if child.Path == path {
			return &child
		}
	}
	return nil
}

/*
func PicByPath(i *Index, path string) *Pic {
  pic := Pic{}
  err := i.dbGetBson("pics", path, &pic)
  if err != nil {
    log.Print(err)
    return nil
  }
  return &pic
}

// Create the given album in the db
func (p Pic) Save(i *Index) {
  log.Printf("Saving pic %s", p.Path)
  err := i.dbSetBson("pics", p.Path, p)
  pics := myDB.Use("Pics")
  // Insert document (document must be map[string]interface{})
  id, err := feeds.Insert()
  if err != nil {
    panic(err)
  }
  p.Id = id
}

// Create the given album in the db
func (a Album) Save(i *Index) {
  log.Printf("Saving album %s", a.Path)
  err := i.dbSetBson("albums", a.Path, a)
  if err != nil {
    panic(err)
  }
}

func AlbumByPath(i *Index, path string) *Album {
  album := Album{}
  err := i.dbGetBson("albums", path, &album)
  if err != nil {
    log.Print(err)
    return nil
  }
  return &album
}

// Props stores global DB props
type Props struct {
  Version int
}*/
