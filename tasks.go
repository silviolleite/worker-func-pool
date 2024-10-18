package workerpool

// Task represents a worker func pool task
type Task struct {
	Name          string
	Fn            any
	Params        []any
	BlockingError bool
}

// Tasks represents the worker func pool tasks
type Tasks struct {
	tasks []Task
}

// NewTasks creates a new worker func pool tasks
func NewTasks() *Tasks {
	return &Tasks{}
}

// AddTask adds a new worker pool task
func (t *Tasks) AddTask(task Task) {
	t.tasks = append(t.tasks, task)
}

// AddTasks adds more than one worker func pool task
func (t *Tasks) AddTasks(task []Task) {
	t.tasks = append(t.tasks, task...)
}

// Tasks returns all worker func pool tasks
func (t *Tasks) Tasks() []Task {
	return t.tasks
}
