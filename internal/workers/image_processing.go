package workers

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"io/ioutil"
	"strings"

	"github.com/nfnt/resize"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/utils/logutil"
)

type imageProcessingWorker struct {
	store    domain.Store
	bucket   string
	filename string
}

func newImageProcessingWorker(s domain.Store, bucket, filename string) *imageProcessingWorker {
	return &imageProcessingWorker{s, bucket, filename}
}

func (i *imageProcessingWorker) Run(ctx context.Context) error {
	r, err := i.store.GetFile(ctx, i.bucket, i.filename)
	if err != nil {
		return fmt.Errorf("[%w] failed to open '%s/%s'", err, i.bucket, i.filename)
	}

	img, _, err := image.Decode(r)
	if err != nil {
		return fmt.Errorf("[%w] failed to decode '%s/%s'", err, i.bucket, i.filename)
	}

	tmp, err := ioutil.TempFile("", "decoded_image")
	if err != nil {
		return fmt.Errorf("[%w] failed to create a temporary file for the decoded image", err)
	}
	defer tmp.Close()

	err = jpeg.Encode(tmp, img, &jpeg.Options{Quality: 100})
	if err != nil {
		return fmt.Errorf("[%w] failed to encode as png file '%s/%s'", err, i.bucket, i.filename)
	}

	logutil.GetDefaultLogger().WithField("tmp name", tmp.Name()).Debug("image encoded as jpg")

	parts := strings.Split(i.filename, ".")
	stats, err := tmp.Stat()
	if err != nil {
		return fmt.Errorf("[%w] failed to get stat for '%s/%s'", err, i.bucket, i.filename)
	}

	_, _ = tmp.Seek(0, 0)

	err = i.store.PutFile(ctx, i.bucket, fmt.Sprintf("%s.jpg", parts[0]), stats.Size(), tmp)
	if err != nil {
		return fmt.Errorf("[%w] failed to save file '%s/%s'", err, i.bucket, i.filename)
	}

	// create the thumbnail
	newImage := resize.Resize(200, 200, img, resize.Lanczos3)

	thumbnailTmp, err := ioutil.TempFile("", "thumbnail_")
	if err != nil {
		return fmt.Errorf("[%w] failed to create tmp file for thumbnail '%s/%s'", err, i.bucket, i.filename)
	}
	defer thumbnailTmp.Close()

	err = jpeg.Encode(thumbnailTmp, newImage, &jpeg.Options{Quality: 100})
	if err != nil {
		return fmt.Errorf("[%w] failed to encode the thumbnail for '%s/%s'", err, i.bucket, i.filename)
	}

	_, _ = thumbnailTmp.Seek(0, 0)

	thumbnailStats, err := thumbnailTmp.Stat()
	if err != nil {
		return fmt.Errorf("[%w] failed to get stat for thumbnail '%s/%s'", err, i.bucket, i.filename)
	}

	err = i.store.PutFile(ctx, i.bucket, fmt.Sprintf("%s_thumbnail.jpg", parts[0]), thumbnailStats.Size(), thumbnailTmp)
	if err != nil {
		return fmt.Errorf("[%w] failed to save thumbnail file '%s/%s'", err, i.bucket, i.filename)
	}

	return nil
}
