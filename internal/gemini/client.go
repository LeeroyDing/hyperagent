package gemini

import (
"context"
"fmt"
"log/slog"
"time"

"github.com/google/generative-ai-go/genai"
"google.golang.org/api/option"
)

type Message struct {
Role    string
Content string
}

type Client struct {
client *genai.Client
model  *genai.GenerativeModel
}

func NewClient(ctx context.Context, apiKey string, modelName string) (*Client, error) {
client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
if err != nil {
return nil, fmt.Errorf("failed to create genai client: %w", err)
}

model := client.GenerativeModel(modelName)
return &Client{
client: client,
model:  model,
}, nil
}

func (c *Client) GenerateContent(ctx context.Context, messages []Message) (string, error) {
cs := c.model.StartChat()

// Map roles and set history except the last message
if len(messages) > 1 {
for _, m := range messages[:len(messages)-1] {
role := "user"
if m.Role == "assistant" || m.Role == "model" {
role = "model"
}
cs.History = append(cs.History, &genai.Content{
Parts: []genai.Part{genai.Text(m.Content)},
Role:  role,
})
}
}

// Last message is the prompt
lastMsg := messages[len(messages)-1]

var lastErr error
for i := 0; i < 3; i++ {
resp, err := cs.SendMessage(ctx, genai.Text(lastMsg.Content))
if err == nil {
if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
return "", fmt.Errorf("no candidates or parts in response")
}
part := resp.Candidates[0].Content.Parts[0]
if text, ok := part.(genai.Text); ok {
return string(text), nil
}
return "", fmt.Errorf("unexpected part type: %T", part)
}
lastErr = err
slog.Warn("Gemini API call failed, retrying...", "attempt", i+1, "error", err)
time.Sleep(time.Duration(1<<i) * time.Second)
}
return "", fmt.Errorf("failed after 3 attempts: %w", lastErr)
}

func (c *Client) EmbedContent(ctx context.Context, text string) ([]float32, error) {
em := c.client.EmbeddingModel("text-embedding-004")
var lastErr error
for i := 0; i < 3; i++ {
resp, err := em.EmbedContent(ctx, genai.Text(text))
if err == nil {
return resp.Embedding.Values, nil
}
lastErr = err
slog.Warn("Gemini Embedding API call failed, retrying...", "attempt", i+1, "error", err)
time.Sleep(time.Duration(1<<i) * time.Second)
}
return nil, fmt.Errorf("failed after 3 attempts: %w", lastErr)
}

func (c *Client) Close() error {
return c.client.Close()
}
