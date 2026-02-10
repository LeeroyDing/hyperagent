package agent

import (
"context"
"errors"
"testing"

"github.com/LeeroyDing/hyperagent/internal/history"
"github.com/stretchr/testify/assert"
)

func TestAgent_Distill(t *testing.T) {
ctx := context.Background()

t.Run("Success", func(t *testing.T) {
h := &MockHistory{
Sessions: map[string][]history.Message{
"s1": make([]history.Message, 6),
},
}
g := &MockGeminiClient{Responses: []string{"summary"}}
m := &MockMemory{}
a := NewAgent(g, nil, m, nil, h, false)

err := a.Distill(ctx, "s1")
assert.NoError(t, err)
assert.Equal(t, "summary", m.Memorized["distill-s1-6"])
})

t.Run("History Load Error", func(t *testing.T) {
h := &MockHistory{LoadError: errors.New("history error")}
a := NewAgent(nil, nil, nil, nil, h, false)
err := a.Distill(ctx, "s1")
assert.Error(t, err)
assert.Contains(t, err.Error(), "history error")
})

t.Run("Not Enough Context", func(t *testing.T) {
h := &MockHistory{
Sessions: map[string][]history.Message{
"s1": make([]history.Message, 3),
},
}
a := NewAgent(nil, nil, nil, nil, h, false)
err := a.Distill(ctx, "s1")
assert.NoError(t, err)
})

t.Run("Gemini Error", func(t *testing.T) {
h := &MockHistory{
Sessions: map[string][]history.Message{
"s1": make([]history.Message, 6),
},
}
g := &MockGeminiClient{GenerateError: errors.New("gemini error")}
a := NewAgent(g, nil, nil, nil, h, false)
err := a.Distill(ctx, "s1")
assert.Error(t, err)
assert.Contains(t, err.Error(), "gemini error")
})

t.Run("Memory Error", func(t *testing.T) {
h := &MockHistory{
Sessions: map[string][]history.Message{
"s1": make([]history.Message, 6),
},
}
g := &MockGeminiClient{Responses: []string{"summary"}}
m := &MockMemory{MemorizeError: errors.New("mem error")}
a := NewAgent(g, nil, m, nil, h, false)
err := a.Distill(ctx, "s1")
assert.Error(t, err)
assert.Contains(t, err.Error(), "mem error")
})
}
