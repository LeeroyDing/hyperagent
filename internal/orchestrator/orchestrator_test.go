package orchestrator

import (
"context"
"errors"
"testing"

"github.com/stretchr/testify/assert"
)

func TestOrchestrator_RunParallel(t *testing.T) {
o := NewOrchestrator()
ctx := context.Background()

tasks := []Task{
{ID: "1", ToolName: "test", ToolArgs: map[string]interface{}{"arg": "val1"}},
{ID: "2", ToolName: "test", ToolArgs: map[string]interface{}{"arg": "val2"}},
{ID: "3", ToolName: "error", ToolArgs: map[string]interface{}{"arg": "val3"}},
}

executeFunc := func(ctx context.Context, task Task) (string, error) {
if task.ToolName == "error" {
return "", errors.New("execution error")
}
return "success: " + task.ID, nil
}

results := o.RunParallel(ctx, tasks, executeFunc)

assert.Len(t, results, 3)
assert.Equal(t, "1", results[0].TaskID)
assert.Equal(t, "success: 1", results[0].Output)
assert.NoError(t, results[0].Error)

assert.Equal(t, "2", results[1].TaskID)
assert.Equal(t, "success: 2", results[1].Output)
assert.NoError(t, results[1].Error)

assert.Equal(t, "3", results[2].TaskID)
assert.Equal(t, "", results[2].Output)
assert.Error(t, results[2].Error)
assert.Equal(t, "execution error", results[2].Error.Error())
}
