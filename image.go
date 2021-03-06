package album

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"

	"code.google.com/p/graphics-go/graphics"
)

var ImageExts = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff"}

// Image service : Reads, saves, scales images and more.
type ImageSvc struct {
}

// CreateScaled creates an image of at most w, h (keeps aspect ratio intact).
// Keeps the original image aspet ratio intact.
// If pad is true, pads with transparency to make the inage size exactly w*h.
func (c ImageSvc) CreateScaled(originalPath, destPath string, w, h int, pad bool) error {
	img, err := c.ReadImage(originalPath)
	if err != nil {
		return err
	}
	img, err = c.ScaledWithin(img, w, h)
	if err != nil {
		return err
	}
	if pad {
		img, err = c.PadImage(img, w, h)
		if err != nil {
			return err
		}
	}
	return c.SaveImage(img, destPath)
}

// IsImage quickly checks if a file looks like an image looking at the file extension.
func (c ImageSvc) IsImage(f os.FileInfo) bool {
	ext := strings.ToLower(filepath.Ext(f.Name()))
	for _, e := range ImageExts {
		if e == ext {
			return true
		}
	}
	return false
}

// ImageSize checks an image dimensions
func (c ImageSvc) ImageSize(filePath string) (x int, y int, err error) {
	img, err := c.ReadImage(filePath)
	if err != nil {
		return 0, 0, err
	}
	return img.Bounds().Dx(), img.Bounds().Dx(), nil
}

// Pad the image with transparency
// Creates a transparent image of width*heigth size with img centered in the center.
// Note that if the image if larger that width, heigth it will get cropped.
func (c ImageSvc) PadImage(img image.Image, width, height int) (i image.Image, err error) {
	padded := image.NewRGBA(image.Rect(0, 0, width, height))
	// DrawMask aligns r.Min in dst with sp in src and mp in mask and then replaces the rectangle r
	// in dst with the result of a Porter-Duff composition. A nil mask is treated as opaque.
	pt := image.Point{
		X: 0,
		Y: 0,
	}
	x := (width - img.Bounds().Dx()) / 2
	y := (height - img.Bounds().Dy()) / 2
	rect := image.Rect(x, y, x+img.Bounds().Dx(), y+img.Bounds().Dy())
	draw.Draw(padded, rect, img, pt, draw.Src)
	return padded, nil
}

// ReadImage eads an image from file
func (c ImageSvc) ReadImage(imgPath string) (img image.Image, err error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return img, err
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return img, err
}

// Rotate rotates the image of path fp and saves it under dest.
func (c ImageSvc) Rotate(fp, dest string, degrees int) (err error) {
	img, err := c.ReadImage(fp)
	if err != nil {
		return err
	}
	img, err = c.Rotated(img, degrees)
	if err != nil {
		return err
	}
	return c.SaveImage(img, dest)
}

// Rotated returns a new image made from rotating ig by n degrees
// Sensible degree values are 90, 180, 270
func (c ImageSvc) Rotated(img image.Image, degrees int) (i image.Image, err error) {
	angle := math.Pi / 180.0 * float64(degrees)
	toImg := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dy(), img.Bounds().Dx()))

	err = graphics.Rotate(toImg, img, &graphics.RotateOptions{Angle: angle})
	if err != nil {
		return toImg, err
	}
	return toImg, err
}

// ReadImageConfig reads the image config from file
func (c ImageSvc) ReadImageConfig(imgPath string) (config image.Config, err error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return config, err
	}
	defer file.Close()
	config, _, err = image.DecodeConfig(file)
	if err != nil {
		return config, err
	}
	return config, nil
}

// SaveImage saves the image usng the proper encoder
// Try to guess the dest image format from it's extension. Uses the encoder default options
// Errors out if the file is not gif, jpeg or png.
func (c ImageSvc) SaveImage(img image.Image, filePath string) error {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".jpg", ".jpeg":
		return c.SaveJpeg(img, filePath, 90)
	case ".png":
		return c.SavePng(img, filePath)
	case ".gif":
		return c.SaveGif(img, filePath, &gif.Options{})
	}
	return fmt.Errorf("Unsupported image file: %s", filePath)
}

// SavePng saves img in PNG format
func (c ImageSvc) SavePng(img image.Image, filePath string) error {
	os.MkdirAll(path.Dir(filePath), 0755)
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, img)
}

// SaveJpeg saves img in JPEG format
func (c ImageSvc) SaveJpeg(img image.Image, filePath string, quality int) error {
	os.MkdirAll(path.Dir(filePath), 0755)
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// write new image to file
	options := jpeg.Options{Quality: quality}
	return jpeg.Encode(out, img, &options)
}

// SaveGif saves img in GF format
func (c ImageSvc) SaveGif(img image.Image, filePath string, options *gif.Options) error {
	os.MkdirAll(path.Dir(filePath), 0755)
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	return gif.Encode(out, img, options)
}

// ScaledWithin returns a scaled version of img.
// It scales img to FIT within width and heigth (whichever is the smallest)
func (c ImageSvc) ScaledWithin(img image.Image, width, height int) (i image.Image, err error) {
	// calculate scaling ratios
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	wr, hr := float32(w)/float32(width), float32(h)/float32(height)
	// using whichever is smaller to fit within width, height
	if wr > hr {
		h = int(float32(h) / wr)
		w = width
	} else {
		w = int(float32(w) / hr)
		h = height
	}
	return c.ScaledImage(img, w, h)
}

// ScaledImage return a new image made of img scaled to width, height
func (c ImageSvc) ScaledImage(img image.Image, width, height int) (i image.Image, err error) {
	// scale
	toImg := image.NewRGBA(image.Rect(0, 0, width, height))
	err = graphics.Scale(toImg, img)
	if err != nil {
		return toImg, err
	}
	return toImg, err
}
