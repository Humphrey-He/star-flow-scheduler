package runtime

import (
	"errors"

	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/model"
)

var ErrQueueFull = errors.New("task queue full")

type TaskQueue struct {
	ch chan *model.Task
}

func NewTaskQueue(size int) *TaskQueue {
	return &TaskQueue{ch: make(chan *model.Task, size)}
}

func (q *TaskQueue) Enqueue(task *model.Task) error {
	select {
	case q.ch <- task:
		return nil
	default:
		return ErrQueueFull
	}
}

func (q *TaskQueue) Channel() <-chan *model.Task {
	return q.ch
}

func (q *TaskQueue) Len() int {
	return len(q.ch)
}
