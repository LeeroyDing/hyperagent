package orchestrator

import (
"context"
"sync"
)

// Task represents a single tool call to be executed.
type Task struct {
ID       string
ToolName string
ToolArgs map[string]interface{}
}

// Result represents the outcome of a task execution.
type Result struct {
TaskID string
Output string
Error  error
}

// Orchestrator manages parallel task execution.
type Orchestrator struct{}

// NewOrchestrator creates a new Orchestrator.
func NewOrchestrator() *Orchestrator {
return &Orchestrator{}
}

// RunParallel executes multiple tasks concurrently.
func (o *Orchestrator) RunParallel(ctx context.Context, tasks []Task, executeFunc func(context.Context, Task) (string, error)) []Result {
var wg sync.WaitGroup
results := make([]Result, len(tasks))

for i, task := range tasks {
wg.Add(1)
go func(idx int, t Task) {
defer wg.Done()
output, err := executeFunc(ctx, t)
results[idx] = Result{
TaskID: t.ID,
Output: output,
Error:  err,
}
}(i, task)
}

wg.Wait()
return results
}
