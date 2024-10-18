package workerpool

import (
	"errors"
	"log"
	"reflect"
)

type taskFunc func() (resp any, err error)

type taskConfig struct {
	name          string
	fn            taskFunc
	blockingError bool
}

// WorkerPool represents the worker pool
// With yours channels and settings fields
type WorkerPool struct {
	taskQueue      chan taskConfig
	resultChan     chan Result
	workerPoolSize int
	workersStarted int
}

// NewWorkerPool creates a new worker pool.
func NewWorkerPool(workerPoolSize int) *WorkerPool {
	return &WorkerPool{
		taskQueue:      make(chan taskConfig),
		resultChan:     make(chan Result),
		workerPoolSize: workerPoolSize,
	}
}

// Start creates and starts the workers.
//
// It uses the workerPoolSize to start the number of workers concurrently.
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerPoolSize; i++ {
		w := worker{id: i, taskQueue: wp.taskQueue, resultChan: wp.resultChan}
		w.start()
		wp.workersStarted++
	}
}

// WorkersRunning returns the number of workers running
func (wp *WorkerPool) WorkersRunning() int {
	return wp.workersStarted
}

// SubmitTask writes the new task to the task queue
//
// parameters:
//   - taskName (the task name).
//   - blockingError (Will return on the result).
//   - fn (The function to execute on the workers).
//   - params (All fn parameters).
func (wp *WorkerPool) SubmitTask(taskName string, blockingError bool, fn any, params ...any) error {
	if fn == nil {
		return errors.New("parameter fn cannot be nil")
	}

	wp.taskQueue <- taskConfig{
		name: taskName,
		fn: func() (resp any, err error) {
			return executeFunc(fn, params...)
		},
		blockingError: blockingError,
	}
	return nil
}

// GetResult returns the worker result
//
// Use it to collect worker results
func (wp *WorkerPool) GetResult() Result {
	return <-wp.resultChan
}

func executeFunc(requestFn any, params ...any) (resp any, err error) {
	fn := reflect.ValueOf(requestFn)
	supportedOutputs := 2
	if fn.Type().NumIn() != len(params) {
		log.Panic("incorrect number of parameters!")
	}
	if fn.Type().NumOut() > supportedOutputs {
		log.Panic("unsupported number of output parameters")
	}

	inputs := makeInputs(params...)
	fnRequest := do(fn, inputs)
	resp, err = fnRequest()

	return resp, err
}

func do(fn reflect.Value, inputs []reflect.Value) func() (any, error) {
	return func() (resp any, err error) {
		values := fn.Call(inputs)

		for _, v := range values {
			if v.Type().Name() == "error" {
				err, _ = v.Interface().(error)
			} else {
				resp = v.Interface()
			}
		}
		return resp, err
	}
}

func makeInputs(params ...any) []reflect.Value {
	inputs := make([]reflect.Value, len(params))
	for k, in := range params {
		inputs[k] = reflect.ValueOf(in)
	}
	return inputs
}
