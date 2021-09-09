package workers

type JobType int

const (
	TransformImageType JobType = iota
)

type Job struct {
	Type JobType
	Data interface{}
}

type JobManager interface {
	Create(j Job) (id string, err error)
	Status(jobID string) (status int, err error)
}
