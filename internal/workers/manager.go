package workers

import (
	"sync"

	"github.com/tupyy/gophoto/internal/domain"
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
