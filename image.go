package album

import (
	"os"
	"path/filepath"
	"strings"
)

var ImageExts = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff"}

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

/*

// Serving Images (With caching + resizing)
type ImageSvc struct {
}

// Serve an Image but scale it on demand
// Size are thumb(200), small(400), large(1024), full
func (c ImageSvc) Serve(prefix, size string, imgPath string) {

  image := filepath.Join(c.IMAGES_FOLDER, size, imgPath)
  if !strings.HasPrefix(image, c.IMAGES_FOLDER) {
    log.Printf("Attempted to read file outside of base imgPath: %s", image)
    return c.NotFound("File not found")
  }

  _, err := os.Stat(image)
  if err != nil {
    if os.IsNotExist(err) {
      CreateScaledImages(imgPath)
      _, err = os.Stat(image)
      if err != nil {
        log.Printf("File not found (%s): %s ", image, err)
        return c.NotFound("File not found")
      }
    } else { // Other unexpected error
      log.Print("Could not access file (%s): %s ", image, err)
      return c.NotFound("File not found")
    }
  }

  file, err := os.Open(image)
  //return c.RenderFile(file, revel.Inline)
}

// copy image from drupal folder and then create scaled versions
func (c ImageSvc) CreateScaledImages(imgPath string) (err error) {
  // create scaled down versions
  dest := filepath.Join(c.IMAGES_FOLDER, "thumb", imgPath)
  resizeImage(dest, thumb, "thumb")
  if err != nil {
    return err
  }
  medium := filepath.Join(c.IMAGES_FOLDER, "medium", imgPath)
  resizeImage(dest, small, "small")
  if err != nil {
    return err
  }
  small := filepath.Join(c.IMAGES_FOLDER, "small", imgPath)
  resizeImage(dest, large, "large")
  if err != nil {
    return err
  }
  return err
}

// Resize the image
func (c ImageSvc) resizeImage(original string, target string, size string) (fname string, err error) {
  file, err := os.Open(original)
  if err != nil {
    return "", err
  }
  img, _, err := image.Decode(file)
  if err != nil {
    return "", err
  }
  file.Close()

  var sz uint = 480
  if size == "medium" {
    sz = 768
  }
  m := resize.Resize(sz, 0, img, resize.Bilinear)

  os.MkdirAll(path.Dir(target), 0755)
  out, err := os.Create(target)
  if err != nil {
    return "", err
  }
  defer out.Close()

  // write new image to file
  log.Printf("Creating scaled down image: %s", target)
  options := jpeg.Options{Quality: 100}
  jpeg.Encode(out, m, &options)

  return fname, nil
}

*/
