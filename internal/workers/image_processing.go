package workers

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"io"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/tupyy/gophoto/utils/logutil"
)

func ProcessImage(r io.Reader, imgWriter, thumbnailWriter io.Writer) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return fmt.Errorf("[%w] failed to decode image", err)
	}

	rotation := getImageRotation(r)
	if rotation != 0 {
		img = imaging.Rotate(img, rotation, color.Gray{})
	}

	err = jpeg.Encode(imgWriter, img, &jpeg.Options{Quality: 100})
	if err != nil {
		return fmt.Errorf("[%w] failed to encode as jpg", err)
	}

	logutil.GetDefaultLogger().Debug("image encoded as jpg")

	// create the thumbnail
	newImage := resize.Resize(200, 0, img, resize.Lanczos3)

	err = jpeg.Encode(thumbnailWriter, newImage, &jpeg.Options{Quality: 100})
	if err != nil {
		return fmt.Errorf("[%w] failed to encode the thumbnail", err)
	}

	return nil
}

func getImageRotation(r io.Reader) float64 {
	x, err := exif.Decode(r)
	var rotation float64 = 0

	if err == nil {
		orientationRaw, err := x.Get("Orientation")

		if err == nil {
			orientation := orientationRaw.String()

			if orientation == "3" {
				rotation = 180
			} else if orientation == "6" {
				rotation = 270
			} else if orientation == "8" {
				rotation = 90
			}
		}
	} else {
		logutil.GetDefaultLogger().WithError(err).Warn("failed to get rotation")
	}

	return rotation
}
