package main

import (
	"fmt"
	"strings"

	workerpool "github.com/silviolleite/worker-func-pool"
)

type MyInt int
type MyString string

func funcA(i int) (MyInt, error) {
	if i == 42 {
		return 0, fmt.Errorf("func A with error")
	}

	return MyInt(i + i), nil
}

func funcB(a, b string) (MyString, error) {
	if a == "error" {
		return "", fmt.Errorf("func B with error")
	}
	return MyString(a + b), nil
}

func funcC(a, b, c string) (string, error) {
	if a == "error" {
		return "", fmt.Errorf("func B with error")
	}

	return strings.Join([]string{a, b, c}, ", "), nil
}

// Data simulates a process data consolidator
type Data struct {
	myInt    MyInt
	myString MyString
	s        string
}

func main() {
	tasks := buildTasks()

	// Creating a worker pool with worker pool size 3
	// but len(tasks.Tasks())) could be used too
	w := workerpool.NewWorkerPool(3)
	w.Start()

	// Submitting the tasks to worker pool
	for _, task := range tasks.Tasks() {
		err := w.SubmitTask(task.Name, task.BlockingError, task.Fn, task.Params...)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	data := Data{}

	for i := 0; i < len(tasks.Tasks()); i++ {
		result := w.GetResult()
		if result.Err != nil {
			if result.BlockingError {
				// Log it as an error and return (In a real case you could return the error)
				fmt.Printf("ERROR: Worker ID: %d, Task: %s, Error: %v\n", result.WorkerID, result.TaskName, result.Err)
				return
			}
			// Just warning because it is not a blocking error
			fmt.Printf("WARN: Worker ID: %d, Task: %s, Error: %v\n", result.WorkerID, result.TaskName, result.Err)
		}

		switch v := result.Data.(type) {
		case MyInt:
			data.myInt = v
			fmt.Printf("Yeah My Int: %v Task: %s, workerID: %v\n", v, result.TaskName, result.WorkerID)
		case MyString:
			data.myString = v
			fmt.Printf("Yeah My String: %v Task: %s, workerID: %v\n", v, result.TaskName, result.WorkerID)
		case string:
			data.s = v
			fmt.Printf("Yeah string: %v Task: %s, workerID: %v\n", v, result.TaskName, result.WorkerID)
		}
	}

	fmt.Printf("\nMy Data: %#v\n", data)

	// output
	// Yeah My Int: 4 Task: A, workerID: 0
	// WARN: Worker ID: 2, Task: B, Error: func B with error
	// Yeah string: a, b, c Task: C, workerID: 1
	//
	// My Data: main.Data{myInt:4, myString:"", s:"a, b, c"}
}

func buildTasks() *workerpool.Tasks {
	// Create a new worker pool Tasks
	tasks := workerpool.NewTasks()

	// Using the method AddTask to add new tasks
	tasks.AddTask(workerpool.Task{
		Name:          "A",
		Fn:            funcA,
		Params:        []any{2}, // Change it to 42 to force error
		BlockingError: false,    // Change it to true if you want to stop the process in case of an error
	})

	tasks.AddTask(workerpool.Task{
		Name:          "B",
		Fn:            funcB,
		Params:        []any{"error", "b"}, // Change first param to any string different of "error" to remove error
		BlockingError: false,               // Change it to true if you want to stop the process in case of an error
	})

	tasks.AddTask(workerpool.Task{
		Name:          "C",
		Fn:            funcC,
		Params:        []any{"a", "b", "c"}, // Change first param to "error" to force error
		BlockingError: false,                // Change it to true if you want to stop the process in case of an error
	})

	return tasks
}
