package parser

import (
"testing"
)

func TestParseLLMResponse(t *testing.T) {
input := "Here is my plan:\n```json\n{\n  \"thoughts\": [\"thinking step 1\"],\n  \"headline\": \"test headline\",\n  \"tool_name\": \"shell\",\n  \"tool_args\": {\"command\": \"ls\"}\n}\n```"

resp, err := ParseLLMResponse(input)
if err != nil {
t.Fatalf("ParseLLMResponse failed: %v", err)
}

if resp.ToolName != "shell" {
t.Errorf("expected tool_name shell, got %s", resp.ToolName)
}

if resp.ToolArgs["command"] != "ls" {
t.Errorf("expected command ls, got %v", resp.ToolArgs["command"])
}
}
