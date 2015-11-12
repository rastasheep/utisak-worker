package worker

import "sync"

const WorkerQueueSize int = 100

type BackgroundWorker struct {
	Queue     chan func()
	waitGroup sync.WaitGroup
}

func NewBackgroundWorker(capacity int) *BackgroundWorker {
	worker := &BackgroundWorker{Queue: make(chan func(), WorkerQueueSize)}

	for i := 0; i < capacity; i++ {
		worker.waitGroup.Add(1)
		go func() {
			for action := range worker.Queue {
				action()
			}
			worker.waitGroup.Done()
		}()
	}
	return worker
}

func (worker *BackgroundWorker) Process() {
	close(worker.Queue)
	worker.waitGroup.Wait()
}
