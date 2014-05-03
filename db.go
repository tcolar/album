// History: Apr 26 14 tcolar Creation

package album

// Album represents a collection of Pics
type Album struct {
	Path         string // Relative path from album root
	Name         string // Pretty name, might or not be same as Path
	Description  string
	HighlightPic string // Album highligh picture relative path
	Ordering     int    // if 0 will order by path
	Hidden       bool   // wether that album/folder should be hidden
	Children     []Album

	// Not serialized along with album structure
	pics  Pics
	dirty bool // Wheter it's "dirty" (pics changed)
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

// Find a picture of the album by path
func (a Album) Pic(path string) *Pic {
	for _, pic := range a.pics {
		if pic.Path == path {
			return &pic
		}
	}
	return nil
}

// Return the pic with the most recent ModTime
// or nil of there are no pics
func (a Album) LatestPic() *Pic {
	var pic *Pic
	for _, p := range a.pics {
		if pic == nil || p.ModTime > pic.ModTime {
			pic = &p
		}
	}
	return pic
}

// Collection of pics with sorting - Sorts the pics by sorting, path
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
	Path        string
	Name        string
	Description string
	Ordering    int  // if 0 will order by path
	Hidden      bool // wether that picture should be hidden / desactivated
	ModTime     int64
}
