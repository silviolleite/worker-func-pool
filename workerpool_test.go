package workerpool_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	workerpool "github.com/silviolleite/worker-func-pool"
)

func TestNewWorkerPool(t *testing.T) {
	w := workerpool.NewWorkerPool(5)
	assert.NotNil(t, w)
}

func TestWorkerPool_Start(t *testing.T) {
	w := workerpool.NewWorkerPool(5)
	assert.NotNil(t, w)
	w.Start()

	assert.Equal(t, 5, w.WorkersRunning())
}

func TestWorkerPool_SubmitTask(t *testing.T) {
	t.Run("FuncNil", func(t *testing.T) {
		w := workerpool.NewWorkerPool(5)
		assert.NotNil(t, w)
		w.Start()
		assert.Equal(t, 5, w.WorkersRunning())

		err := w.SubmitTask("", false, nil)
		wantErr := fmt.Errorf("parameter fn cannot be nil")
		assert.Equal(t, wantErr, err)
	})

	t.Run("SuccessSubmit", func(t *testing.T) {
		w := workerpool.NewWorkerPool(5)
		assert.NotNil(t, w)
		w.Start()
		assert.Equal(t, 5, w.WorkersRunning())

		err := w.SubmitTask("", false, stubFuncA, []any{"a", "b"}...)
		assert.Nil(t, err)
	})
}

func TestWorkerPool_GetResult(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		type Data struct {
			myInt    int
			myString string
		}

		tests := []struct {
			name     string
			tasks    []workerpool.Task
			wantData Data
			wantErr  error
		}{
			{
				name: "SuccessWithOneTask",
				tasks: []workerpool.Task{
					{Name: "a", BlockingError: false, Fn: stubFuncA, Params: []any{"a", "b"}},
				},
				wantData: Data{
					myInt:    0,
					myString: "ab",
				},
				wantErr: nil,
			},
			{
				name: "SuccessWithOneTask",
				tasks: []workerpool.Task{
					{Name: "a", BlockingError: false, Fn: stubFuncA, Params: []any{"a", "b"}},
					{Name: "b", BlockingError: false, Fn: stubFuncB, Params: []any{2}},
				},
				wantData: Data{
					myInt:    2,
					myString: "ab",
				},
				wantErr: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := workerpool.NewWorkerPool(len(tt.tasks))
				w.Start()
				for _, task := range tt.tasks {
					_ = w.SubmitTask(task.Name, task.BlockingError, task.Fn, task.Params...)
				}

				data := Data{}
				for i := 0; i < len(tt.tasks); i++ {
					result := w.GetResult()
					assert.NotNil(t, result)
					assert.NoError(t, result.Err)

					switch v := result.Data.(type) {
					case int:
						data.myInt = v
					case string:
						data.myString = v
					}
				}
				assert.Equal(t, tt.wantData, data)
			})
		}
	})

	t.Run("Error", func(t *testing.T) {
		tests := []struct {
			name              string
			tasks             []workerpool.Task
			wantErr           error
			wantBlockingError bool
		}{
			{
				name: "ErrorBlockingErrorFalse",
				tasks: []workerpool.Task{
					{Name: "a", BlockingError: false, Fn: stubFuncA, Params: []any{"error", "b"}},
				},
				wantErr:           fmt.Errorf("func A with error"),
				wantBlockingError: false,
			},
			{
				name: "ErrorBlockingErrorTrue",
				tasks: []workerpool.Task{
					{Name: "a", BlockingError: true, Fn: stubFuncA, Params: []any{"error", "b"}},
				},
				wantErr:           fmt.Errorf("func A with error"),
				wantBlockingError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := workerpool.NewWorkerPool(len(tt.tasks))
				w.Start()
				for _, task := range tt.tasks {
					_ = w.SubmitTask(task.Name, task.BlockingError, task.Fn, task.Params...)
				}

				for i := 0; i < len(tt.tasks); i++ {
					result := w.GetResult()
					assert.NotNil(t, result.Err)
					assert.ErrorContains(t, result.Err, tt.wantErr.Error())
					assert.Equal(t, tt.wantBlockingError, result.BlockingError)
				}
			})
		}
	})
}
