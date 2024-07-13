package internal

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type LockableTasks struct {
	tasksMap            []bool
	editLock            sync.Mutex
	concurrentTaskCount int
}

type Task struct {
	Id       int
	tasksMap *LockableTasks
}

func NewLockableTasks(tasksCount int) LockableTasks {
	log.Info().Msg("initialized the LockableTasks system")
	return LockableTasks{
		tasksMap:            make([]bool, tasksCount),
		editLock:            sync.Mutex{},
		concurrentTaskCount: tasksCount,
	}
}
func (r *LockableTasks) AssignSessionId(ctx context.Context) (*Task, error) {
	ticker := time.NewTicker(time.Millisecond)
	counter := 0
	log.Debug().Msg("starting allocating a task number")
	for {
		select {
		case <-ticker.C:
			log.Debug().Int("try", counter).Int("task-id", counter%r.concurrentTaskCount).Msg("trying ID")

			if counter > 3*r.concurrentTaskCount {
				// We do not allow any more searching to avoid infinite loops
				return nil, fmt.Errorf("task allocation timeout reached")
			}

			if r.tasksMap[counter%r.concurrentTaskCount] {
				counter += 1
				continue
			} else {

				// Using a nested function to use `defer` sooner
				func() {
					r.editLock.Lock()
					defer r.editLock.Unlock()

					r.tasksMap[counter%r.concurrentTaskCount] = true
				}()

				task := Task{
					tasksMap: r,
					Id:       counter % r.concurrentTaskCount,
				}
				log.Debug().Int("try", counter).Int("task-id", counter%r.concurrentTaskCount).Msg("finished allocating task id")

				return &task, nil
			}

		case <-ctx.Done():
			log.Debug().Msg("task allocation cancelled by context")
			return nil, ctx.Err()
		}
	}
}

func (r Task) Release() {
	log.Debug().Int("task-id", r.Id).Msg("releasing task")
	r.tasksMap.editLock.Lock()
	defer r.tasksMap.editLock.Unlock()
	r.tasksMap.tasksMap[r.Id] = false
}
