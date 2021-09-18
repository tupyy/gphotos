package workers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type State int

const (
	Idle State = iota
	Running
	Done
)

type Worker interface {
	Run(ctx context.Context) error
}

type Job struct {
	// ID of the job
	ID uuid.UUID
	// State holds the current state of the job
	State State
	// Err is the last error from workers
	Err error
	// Workers holds the list of workers
	workers []Worker
}

func newJob(workers ...Worker) *Job {
	id := uuid.New()
	return &Job{ID: id, State: Idle, workers: workers}
}

func (j *Job) Run(ctx context.Context) error {
	j.State = Running
	defer func() {
		j.State = Done
	}()

	for _, w := range j.workers {
		err := w.Run(ctx)
		if err != nil {
			return fmt.Errorf("%w job_id: %s", err, j.ID.String())
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("job_id: %s. context cancelled: %w", j.ID.String(), ctx.Err())
		default:
		}
	}

	return nil
}
