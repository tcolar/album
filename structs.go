// History: Apr 26 14 tcolar Creation

package album

// Album represents a collection of Pics
type Album struct {
	Id           string // Unique storer key
	ParentId     string // key of parent album
	Path         string // Relative path from album root
	Name         string // Pretty name, defualts to folder name
	Description  string
	HighlightPic string // Album highlight picture key.
	Ordering     int    // if equal secondary ordering is by name
	Hidden       bool   // whether that album/folder should be hidden
}

func (a Album) Pics() []Pic {
	return []Pic{} // TODO
}

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
	Width       int // width of origina in pixel
	Height      int // width of original in pixels
}
