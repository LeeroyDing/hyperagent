package history

import (
"encoding/json"
"fmt"
"os"
"path/filepath"
"sort"
"time"
)

type Message struct {
Role    string    `json:"role"` // "user", "assistant", "system", "tool"
Content string    `json:"content"`
Time    time.Time `json:"time"`
}

type Session struct {
ID        string    `json:"id"` 
Name      string    `json:"name"` 
UpdatedAt time.Time `json:"updated_at"` 
Messages  []Message `json:"messages,omitempty"` 
}

type HistoryManager struct {
StorageDir string
}

func GetDefaultHistoryDir() string {
home, _ := os.UserHomeDir()
return filepath.Join(home, ".hyperagent", "history")
}

func NewHistoryManager(storageDir string) (*HistoryManager, error) {
if storageDir == "" {
storageDir = GetDefaultHistoryDir()
}
if err := os.MkdirAll(storageDir, 0755); err != nil {
return nil, fmt.Errorf("failed to create storage directory: %w", err)
}
return &HistoryManager{StorageDir: storageDir}, nil
}

func (h *HistoryManager) GetSessionPath(sessionID string) string {
return filepath.Join(h.StorageDir, sessionID+".jsonl")
}

func (h *HistoryManager) GetMetadataPath(sessionID string) string {
return filepath.Join(h.StorageDir, sessionID+".meta.json")
}

func (h *HistoryManager) AddMessage(sessionID, role, content string) error {
msg := Message{
Role:    role,
Content: content,
Time:    time.Now(),
}

path := h.GetSessionPath(sessionID)
f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
return fmt.Errorf("failed to open history file: %w", err)
}
defer f.Close()

data, err := json.Marshal(msg)
if err != nil {
return fmt.Errorf("failed to marshal message: %w", err)
}

if _, err := f.Write(append(data, '\n')); err != nil {
return fmt.Errorf("failed to write to history file: %w", err)
}

return nil
}

func (h *HistoryManager) SetSessionName(sessionID, name string) error {
path := h.GetMetadataPath(sessionID)
meta := map[string]string{"name": name}
data, err := json.Marshal(meta)
if err != nil {
return err
}
return os.WriteFile(path, data, 0644)
}

func (h *HistoryManager) GetSessionName(sessionID string) string {
path := h.GetMetadataPath(sessionID)
data, err := os.ReadFile(path)
if err != nil {
return "New Conversation"
}
var meta map[string]string
if err := json.Unmarshal(data, &meta); err != nil {
return "New Conversation"
}
return meta["name"]
}

func (h *HistoryManager) LoadHistory(sessionID string) ([]Message, error) {
path := h.GetSessionPath(sessionID)
if _, err := os.Stat(path); os.IsNotExist(err) {
return []Message{}, nil
}

f, err := os.Open(path)
if err != nil {
return nil, fmt.Errorf("failed to open history file: %w", err)
}
defer f.Close()

messages := make([]Message, 0)
decoder := json.NewDecoder(f)
for decoder.More() {
var msg Message
if err := decoder.Decode(&msg); err != nil {
return nil, fmt.Errorf("failed to decode message: %w", err)
}
messages = append(messages, msg)
}

return messages, nil
}

func (h *HistoryManager) ListSessions() ([]Session, error) {
files, err := os.ReadDir(h.StorageDir)
if err != nil {
return nil, fmt.Errorf("failed to read storage directory: %w", err)
}

sessions := make([]Session, 0)
for _, f := range files {
if f.IsDir() || filepath.Ext(f.Name()) != ".jsonl" {
continue
}

info, err := f.Info()
if err != nil {
continue
}

id := f.Name()[:len(f.Name())-6]
sessions = append(sessions, Session{
ID:        id,
Name:      h.GetSessionName(id),
UpdatedAt: info.ModTime(),
})
}

sort.Slice(sessions, func(i, j int) bool {
return sessions[i].UpdatedAt.After(sessions[j].UpdatedAt)
})

return sessions, nil
}
