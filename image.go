package album

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"strings"

	"code.google.com/p/graphics-go/graphics"
)

var ImageExts = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff"}

// Serving Images (With caching + resizing)
type ImageSvc struct {
}

/*
// TODO: ResizeImageWithin with padding ?
// would create an image of exactly width*height padded with transparency.
// obviously would not work wth jpg though
func (c ImageSvc) ResizeImageWithin(original, target string, width, height int) error {
}
*/

// ResizeImage resizes the image original image and saves it as target
// It makes the source image FIT within width and heigth (whihchever is the smallest)
// Try to keep the original image format from it's extension. Uses the encoder default options
// Errors out if the file is not gif, jpeg or png.
func (c ImageSvc) ResizeImageWithin(original, target string, width, height int) error {
	file, err := os.Open(original)
	if err != nil {
		return err
	}
	defer file.Close()
	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return err
	}
	// calculate scaling ratios
	w, h := width, height
	wr, hr := float32(config.Width)/float32(width), float32(config.Height)/float32(height)
	// using whichever is smaller to fit within width, height
	if wr > hr {
		h = int(float32(h) / hr)
	} else {
		w = int(float32(w) / wr)
	}
	return c.ResizeImage(original, target, w, h)
}

// ResizeImage resizes the image original image and saves it as target
// Try to guess the dest image format from it's extension. Uses the encoder default options
// Errors out if the file is not gif, jpeg or png.
func (c ImageSvc) ResizeImage(original, target string, width, height int) error {
	ext := strings.ToLower(filepath.Ext(target))
	switch ext {
	case ".jpg", ".jpeg":
		return c.ResizeImageJpeg(original, target, width, height, 90)
	case ".png":
		return c.ResizeImagePng(original, target, width, height)
	case ".gif":
		return c.ResizeImageGif(original, target, width, height, &gif.Options{})
	}
	return fmt.Errorf("Unsupported image file: %s", original)
}

// ResizeImageGif resizes the original image and saves it as target in Png format
func (c ImageSvc) ResizeImagePng(original, target string, width, height int) error {
	img, err := c.resizeImage(original, target, width, height)
	if err != nil {
		return err
	}
	os.MkdirAll(path.Dir(target), 0755)
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, img)
}

// ResizeImageGif resizes the original image and saves it as target in Jpeg format
func (c ImageSvc) ResizeImageJpeg(original, target string, width, height int, quality int) error {
	img, err := c.resizeImage(original, target, width, height)
	if err != nil {
		return err
	}
	os.MkdirAll(path.Dir(target), 0755)
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()

	// write new image to file
	options := jpeg.Options{Quality: quality}
	return jpeg.Encode(out, img, &options)
}

// ResizeImageGif resizes the original image and saves it as target in Gif format
func (c ImageSvc) ResizeImageGif(original, target string, width, height int, options *gif.Options) error {
	img, err := c.resizeImage(original, target, width, height)
	if err != nil {
		return err
	}
	os.MkdirAll(path.Dir(target), 0755)
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()

	return gif.Encode(out, img, options)
}

func (c ImageSvc) resizeImage(original, target string, width, height int) (img image.Image, err error) {
	file, err := os.Open(original)
	if err != nil {
		return img, err
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	if err != nil {
		return img, err
	}
	// scale
	toImg := image.NewRGBA64(image.Rect(0, 0, width, height))
	err = graphics.Scale(toImg, img)
	if err != nil {
		return toImg, err
	}
	return toImg, err
}

// IsImage quickly checks if a file looks like an image looking at the file extension.
func IsImage(f os.FileInfo) bool {
	ext := strings.ToLower(filepath.Ext(f.Name()))
	for _, e := range ImageExts {
		if e == ext {
			return true
		}
	}
	return false
}
