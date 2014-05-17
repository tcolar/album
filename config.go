// History: Apr 26 14 tcolar Creation

package album

// AlbumConfig: Album server configuration
// TODO: At some point will need a proper config file or such
type AlbumConfig struct {
	AlbumDir      string
	Port          int
	AdminPassword string

	// List of different "sizes"(key) we will scale images too
	MediaSizes map[string]MediaSizing

	// Do we Keep the original or remove it after scaling is done ?
	// TODO : RemoveOriginal bool
}

// Defines image scaling in relation to viewport size (HTML media queries)
type MediaSizing struct {
	// MinMediaWidth specify the minimum media width to use this sizing
	// For example 1024 would ransalet to -> media="(min-width: 1024px)
	MinMediaWidth int
	// What width (at most) should the image be scaled to for this.
	MaxScaleWidth int
	// What height (at most) should the image be scaled to for this.
	MaxScaleHeight int
	// Whether to pad the image with tranaparency to make it EXACTLY the given size
	PadOnScale bool
}

// Default sizes for photo album
var StdAlbumSizes map[string]MediaSizing = map[string]MediaSizing{
	"thumb":  MediaSizing{0, 200, 200, true},
	"small":  MediaSizing{0, 600, 900, false},
	"medium": MediaSizing{640, 1000, 1400, false},
	"large":  MediaSizing{1280, 1440, 1440, false},
}
