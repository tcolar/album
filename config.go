// History: Apr 26 14 tcolar Creation

package album

// AlbumConfig: Album server configuration
// TODO: At some point will need a proper config file or such
type AlbumConfig struct {
	AlbumDir      string
	Port          int
	AdminPassword string

	// Thumbnails of at most "thumbsize" will go in /_scaled/thumb/
	ThumbSize int

	// See http://foundation.zurb.com/docs/media-queries.html for small, medium, large specs
	// Scaled small image will be served to small devices (~phones) /_scaled/small/
	SmallSize MediaSizing
	// Scaled medium inages to be served to midsize devices (~tablets) /_scaled/medium/
	MedSize MediaSizing
	// Scaled large inages to be served to large devices (~PC)
	// Will either be the original if ResizeOriginal is true or a new image under /_scaled/large/
	LargeSize MediaSizing
	// Do we Resize and REPLACE the originals to be the "LargeSize" or leave them intact ?
	// In which case a new file will be created for large under /_scaled/large/
	ResizeOriginal bool
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
}

/*

/*
<picture>
  <source srcset="examples/images/large.jpg" media="(min-width: 1024px)">
  <source srcset="examples/images/medium.jpg" media="(min-width: 640px)">
  <source srcset="examples/images/small.jpg">
  <img srcset="examples/images/medium.jpg" alt="A giant stone face at The Bayon temple in Angkor Thom, Cambodia">
</picture>
*/
