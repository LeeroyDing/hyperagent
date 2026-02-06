package executor

import (
"bufio"
"fmt"
"io"
"os"
"os/exec"
"sync"

"github.com/creack/pty"
)

// ShellSession represents a persistent PTY session
type ShellSession struct {
ID     string
Cmd    *exec.Cmd
Pty    *os.File
Reader *bufio.Reader
mu     sync.Mutex
}

// NewShellSession spawns a new persistent shell
func NewShellSession(id string) (*ShellSession, error) {
c := exec.Command("bash", "--noprofile", "--norc")
// Set a custom prompt to make detection easier
c.Env = append(os.Environ(), "PS1=HYPERAGENT_PROMPT> ")

f, err := pty.Start(c)
if err != nil {
return nil, err
}

s := &ShellSession{
ID:     id,
Cmd:    c,
Pty:    f,
Reader: bufio.NewReader(f),
}

// Wait for the initial prompt
_, err = s.readUntilPrompt()
if err != nil {
s.Close()
return nil, fmt.Errorf("failed to initialize shell: %v", err)
}

return s, nil
}

func (s *ShellSession) Execute(command string) (string, error) {
s.mu.Lock()
defer s.mu.Unlock()

// Write command to PTY
_, err := fmt.Fprintln(s.Pty, command)
if err != nil {
return "", err
}

return s.readUntilPrompt()
}

func (s *ShellSession) readUntilPrompt() (string, error) {
var output []byte
prompt := "HYPERAGENT_PROMPT> "

for {
line, err := s.Reader.ReadBytes('\n')
if err != nil {
if err == io.EOF {
break
}
return string(output), err
}
output = append(output, line...)

// Check if the last part of the output contains the prompt
// Note: This is a simple implementation. A more robust one would handle
// the prompt appearing without a newline (e.g. after a command finishes).
if s.Reader.Buffered() == 0 && containsPrompt(output, prompt) {
break
}
}

result := string(output)
// Clean up the output by removing the echoed command and the prompt
// (Simplified for prototype)
return result, nil
}

func containsPrompt(output []byte, prompt string) bool {
return len(output) >= len(prompt) && string(output[len(output)-len(prompt):]) == prompt
}

func (s *ShellSession) Close() error {
s.Pty.Close()
return s.Cmd.Process.Kill()
}

// SessionManager manages multiple shell sessions
type SessionManager struct {
sessions map[string]*ShellSession
mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
return &SessionManager{
sessions: make(map[string]*ShellSession),
}
}

func (m *SessionManager) GetOrCreate(id string) (*ShellSession, error) {
m.mu.Lock()
defer m.mu.Unlock()

if s, ok := m.sessions[id]; ok {
return s, nil
}

s, err := NewShellSession(id)
if err != nil {
return nil, err
}
m.sessions[id] = s
return s, nil
}

func (m *SessionManager) Cleanup() {
m.mu.Lock()
defer m.mu.Unlock()
for _, s := range m.sessions {
s.Close()
}
}
