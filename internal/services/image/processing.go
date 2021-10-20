package image

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	"io"

	"github.com/disintegration/imaging"
	"github.com/tupyy/gophoto/utils/logutil"
)

// Process encode the image as jpg and create a thumbnail.
// It also tries to rotate the image if the "Orientation" is set in the Exif.
func Process(r io.Reader, imgWriter, thumbnailWriter io.Writer) error {
	img, err := imaging.Decode(r, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("%w failed to decode image", err)
	}

	err = imaging.Encode(imgWriter, img, imaging.JPEG, imaging.JPEGQuality(80))
	if err != nil {
		return fmt.Errorf("failed to encode as jpg: %v", err)
	}

	logutil.GetDefaultLogger().Debug("image encoded as jpg")

	// create the thumbnail
	thumbnail := imaging.Resize(img, 0, 150, imaging.Lanczos)

	err = imaging.Encode(thumbnailWriter, thumbnail, imaging.JPEG, imaging.JPEGQuality(100))
	if err != nil {
		return fmt.Errorf("failed to encode the thumbnail: %v", err)
	}

	logutil.GetDefaultLogger().Debug("thumbnail image created")

	return nil
}
