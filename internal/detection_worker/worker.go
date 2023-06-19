package detection_worker

import "github.com/go-co-op/gocron"

// Worker -.
type Worker struct {
	scheduler *gocron.Scheduler
}

// NewWorker -.
func New() *Worker {
	return &Worker{}
}
