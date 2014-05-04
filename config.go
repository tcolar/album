// History: Apr 26 14 tcolar Creation

package album

// AlbumConfig: configuration
// TODO: Load from file and/or command line flags
type AlbumConfig struct {
	AlbumDir string
	DbDir    string
	// TODO: make gallery server optional ?
	// TODO: make api server optional ?
}
