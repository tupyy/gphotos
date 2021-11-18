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
func Process(r io.ReadSeeker, imgWriter io.Writer) error {
	if _, err := r.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to decode iamge: %v", err)
	}

	img, err := imaging.Decode(r, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	err = imaging.Encode(imgWriter, img, imaging.JPEG, imaging.JPEGQuality(90))
	if err != nil {
		return fmt.Errorf("failed to encode as jpg: %v", err)
	}

	logutil.GetDefaultLogger().Debug("image encoded as jpg")

	return nil
}

func CreateThumbnail(r io.ReadSeeker, w io.Writer) error {
	if _, err := r.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to create thumbnail: %v", err)
	}

	img, err := imaging.Decode(r, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	// create the thumbnail
	thumbnail := imaging.Resize(img, 0, 150, imaging.Lanczos)

	err = imaging.Encode(w, thumbnail, imaging.JPEG, imaging.JPEGQuality(100))
	if err != nil {
		return fmt.Errorf("failed to encode the thumbnail: %v", err)
	}

	logutil.GetDefaultLogger().Debug("thumbnail image created")

	return nil
}
