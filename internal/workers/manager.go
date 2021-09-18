package workers

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/utils/logutil"
)

var (
	jobManager    *JobManager
	createJobOnce sync.Once
)

type JobManager struct {
	// MaxConcurrentJobs maximum number of concurrent jobs
	maxConcurrentJobs int
	// jobs holds the list of jobs
	jobs []*Job
	// store is the interface to domain.Store
	store domain.Store
}

func NewJobManager(maxConcurrentJobs int, store domain.Store) *JobManager {
	createJobOnce.Do(func() {
		jobManager = &JobManager{maxConcurrentJobs: maxConcurrentJobs, store: store}
	})

	return jobManager
}

func (jm *JobManager) NewImageProcessingJob(bucket, filename, dstBucket string) string {
	// create processing worker
	processingWorker := newImageProcessingWorker(jm.store, bucket, filename)

	parts := strings.Split(filename, ".")
	copyImageWorker := newImageCopyWorker(jm.store, bucket, fmt.Sprintf("%s.jpg", parts[0]), dstBucket)
	copyThumbnailWorker := newImageCopyWorker(jm.store, bucket, fmt.Sprintf("%s_thumbnail.jpg", parts[0]), dstBucket)

	j := newJob(processingWorker, copyImageWorker, copyThumbnailWorker)
	jm.jobs = append(jm.jobs, j)

	// temporary
	err := j.Run(context.Background())
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Errorf("failed to process image '%s/%s'", bucket, filename)
	}

	return j.ID.String()
}
