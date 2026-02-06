package token

import (
"fmt"
"github.com/pkoukk/tiktoken-go"
)

// TokenManager handles token counting and pruning.
type TokenManager struct {
encoding *tiktoken.Tiktoken
}

// NewTokenManager creates a new TokenManager.
func NewTokenManager(model string) (*TokenManager, error) {
// Gemini models don't have a direct tiktoken encoding, but cl100k_base is a good approximation for many LLMs.
// For more accuracy, one would use the Gemini API's countTokens method.
enc, err := tiktoken.GetEncoding("cl100k_base")
if err != nil {
return nil, fmt.Errorf("failed to get encoding: %w", err)
}
return &TokenManager{encoding: enc}, nil
}

// CountTokens returns the number of tokens in a string.
func (tm *TokenManager) CountTokens(text string) int {
tokens := tm.encoding.Encode(text, nil, nil)
return len(tokens)
}

// PruneHistory prunes the history messages to fit within the token limit.
// This is a placeholder for more complex pruning logic.
func (tm *TokenManager) PruneHistory(messages []string, maxTokens int) []string {
totalTokens := 0
var pruned []string
for i := len(messages) - 1; i >= 0; i-- {
tokens := tm.CountTokens(messages[i])
if totalTokens+tokens > maxTokens {
break
}
totalTokens += tokens
pruned = append([]string{messages[i]}, pruned...)
}
return pruned
}
