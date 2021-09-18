package workers

import (
	"context"
	"fmt"

	"github.com/tupyy/gophoto/internal/domain"
)

type ImageCopyWorker struct {
	store       domain.Store
	srcBucket   string
	srcFilename string
	dstBucket   string
}

func newImageCopyWorker(s domain.Store, srcBucket, srcFilename, dstBucket string) *ImageCopyWorker {
	return &ImageCopyWorker{s, srcBucket, srcFilename, dstBucket}
}

func (i *ImageCopyWorker) Run(ctx context.Context) error {
	r, err := i.store.GetFile(ctx, i.srcBucket, i.srcFilename)
	if err != nil {
		return fmt.Errorf("%w failed to open file '%s/%s'", err, i.srcBucket, i.srcFilename)
	}

	err = i.store.PutFile(ctx, i.dstBucket, i.srcFilename, -1, r)
	if err != nil {
		return fmt.Errorf("%w failed to copy '%s/%s' to '%s'", err, i.srcBucket, i.srcFilename, i.dstBucket)
	}

	// remove original file
	err = i.store.DeleteFile(ctx, i.srcBucket, i.srcFilename)
	if err != nil {
		return fmt.Errorf("%w failed to delete original file '%s/%s' after move", err, i.srcBucket, i.srcFilename)
	}

	return nil
}
