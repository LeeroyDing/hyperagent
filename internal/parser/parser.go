package parser

import (
"encoding/json"
"fmt"
"regexp"
"strings"

"github.com/tidwall/gjson"
)

// Response represents the structured response from the LLM.
type Response struct {
Thoughts []string               `json:"thoughts"` 
Headline string                 `json:"headline"` 
ToolName string                 `json:"tool_name"` 
ToolArgs map[string]interface{} `json:"tool_args"` 
}

// ParseLLMResponse extracts and parses JSON from the LLM's output.
func ParseLLMResponse(input string) (*Response, error) {
// Try to find JSON in markdown code blocks
re := regexp.MustCompile("```(?:json)?\\s*([\\s\\S]*?)```")
matches := re.FindStringSubmatch(input)

jsonStr := input
if len(matches) > 1 {
jsonStr = matches[1]
}

jsonStr = strings.TrimSpace(jsonStr)

// Use gjson to validate it's at least a JSON object
if !gjson.Valid(jsonStr) {
return nil, fmt.Errorf("invalid JSON: %s", jsonStr)
}

var resp Response
err := json.Unmarshal([]byte(jsonStr), &resp)
if err != nil {
return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
}

return &resp, nil
}
