package workerpool_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	workerpool "github.com/silviolleite/worker-func-pool"
)

func stubFuncA(a, b string) (string, error) {
	if a == "error" {
		return "", fmt.Errorf("func A with error")
	}

	return a + b, nil
}

func stubFuncB(n int) (int, error) {
	if n == 42 {
		return 0, fmt.Errorf("func B with error")
	}

	return n, nil
}

func TestNewTasks(t *testing.T) {
	tasks := workerpool.NewTasks()
	assert.NotNil(t, tasks)
}

func TestTasks_AddTask(t *testing.T) {
	tasks := workerpool.NewTasks()
	tasks.AddTask(workerpool.Task{
		Name:          "A",
		Fn:            stubFuncA,
		Params:        []any{"test", "ing"},
		BlockingError: false,
	})
	assert.Len(t, tasks.Tasks(), 1)
}

func TestTasks_Run(t *testing.T) {
	tests := []struct {
		name    string
		tasks   []workerpool.Task
		wantLen int
	}{
		{
			name:    "WithTasksNil",
			tasks:   nil,
			wantLen: 0,
		},
		{
			name:    "WithOneTask",
			tasks:   []workerpool.Task{{Name: "A"}},
			wantLen: 1,
		},
		{
			name:    "WithMoreThanOneTask",
			tasks:   []workerpool.Task{{Name: "A"}, {Name: "B"}},
			wantLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks := workerpool.NewTasks()
			tasks.AddTasks(tt.tasks)
			assert.Len(t, tasks.Tasks(), tt.wantLen)
		})
	}
}
