package album

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

var ImageExts = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff"}

// Serving Images (With caching + resizing)
type ImageSvc struct {
}

// ResizeImageGif resizes the original image and saves it as target in Png format
func (c ImageSvc) ResizeImagePng(original, target string, width, height uint) error {
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
func (c ImageSvc) ResizeImageJpeg(original, target string, width, height uint, quality int) error {
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
func (c ImageSvc) ResizeImageGif(original, target string, width, height uint, options *gif.Options) error {
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

func (c ImageSvc) resizeImage(original, target string, width, height uint) (img image.Image, err error) {
	file, err := os.Open(original)
	if err != nil {
		return img, err
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	if err != nil {
		return img, err
	}

	// TODO: let user pick interpolation function ?
	return resize.Resize(width, height, img, resize.Bilinear), nil
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
