package workerpool

// Result represents the worker processed data
type Result struct {
	WorkerID      int
	TaskName      string
	Data          any
	Err           error
	BlockingError bool
}

type worker struct {
	id         int
	taskQueue  <-chan taskConfig
	resultChan chan<- Result
}

func (w *worker) start() {
	go func() {
		mustStop := false
		for t := range w.taskQueue {
			data, err := t.fn()
			if err != nil && t.blockingError {
				mustStop = true
			}

			w.resultChan <- Result{
				WorkerID:      w.id,
				TaskName:      t.name,
				BlockingError: t.blockingError,
				Err:           err,
				Data:          data,
			}
		}

		if mustStop {
			w.stop()
		}
	}()
}

func (w *worker) stop() {
	close(w.resultChan)
}
