package parser

import (
"testing"

"github.com/stretchr/testify/assert"
)

func TestParseLLMResponse_ExtraCoverage(t *testing.T) {
t.Run("Direct JSON input", func(t *testing.T) {
input := `{"thoughts": ["test"], "headline": "test", "tool_name": "test", "tool_args": {}}`
resp, err := ParseLLMResponse(input)
assert.NoError(t, err)
assert.NotNil(t, resp)
assert.Equal(t, "test", resp.ToolName)
})

t.Run("Invalid JSON", func(t *testing.T) {
input := "not a json"
_, err := ParseLLMResponse(input)
assert.Error(t, err)
assert.Contains(t, err.Error(), "invalid JSON")
})

t.Run("Unmarshal error", func(t *testing.T) {
// thoughts is expected to be []string, providing a number should trigger unmarshal error
input := `{"thoughts": 123}`
_, err := ParseLLMResponse(input)
assert.Error(t, err)
assert.Contains(t, err.Error(), "failed to unmarshal JSON")
})

t.Run("Markdown with no JSON tag", func(t *testing.T) {
input := "```\n{\"headline\": \"no-tag\"}\n```"
resp, err := ParseLLMResponse(input)
assert.NoError(t, err)
assert.Equal(t, "no-tag", resp.Headline)
})
}
