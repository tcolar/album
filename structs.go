// History: Apr 26 14 tcolar Creation

package album

// Album represents a collection of Pics
type Album struct {
	Id             string // Unique storer key
	ParentId       string // key of parent album
	Path           string // Relative path from album root
	Name           string // Pretty name, defualts to folder name
	Description    string
	HighlightPic   string // Album highlight picture key.
	Ordering       int    // if equal secondary ordering is by name
	Hidden         bool   // whether that album/folder should be hidden
	LastPicModTime int64  // timestamp of most recent picture in the album
}

func (a Album) Pics() []Pic {
	return []Pic{} // TODO
}

/*
// Album return a child album album by path, or nil if none found.
func (a Album) Child(path string) *Album {
  for _, child := range a.Children {
    if child.Path == path {
      return &child
    }
  }
  return nil
}

// Pic finds a picture of the album by path
func (a Album) Pic(path string) *Pic {
  for _, pic := range a.pics {
    if pic.Path == path {
      return &pic
    }
  }
  return nil
}

// LatestPic returns the pic with the most recent ModTime
// or nil of there are no pics
func (a Album) LatestPic() *Pic {
  var pic *Pic
  for i, _ := range a.pics {
    p := a.pics[i]
    if pic == nil || p.ModTime > pic.ModTime {
      pic = &p
    }
  }
  return pic
}*/

// Collection of picture with sorting. Sorts the pics by sorting, path
type Pics []Pic

func (p Pics) Len() int      { return len(p) }
func (p Pics) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p Pics) Less(i, j int) bool {
	if p[i].Ordering == p[j].Ordering {
		return p[i].Name < p[j].Name
	}
	return p[i].Ordering < p[j].Ordering
}

// Pic representents an individual picture / file
type Pic struct {
	Id          string // Unique storer key
	AlbumId     string // key of album
	Path        string
	Name        string
	Description string
	Ordering    int  // if 0 will order by path
	Hidden      bool // whether that picture should be hidden / desactivated
	ModTime     int64
}
