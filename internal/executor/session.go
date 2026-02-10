package executor

import (
"bufio"
"bytes"
"context"
"fmt"
"os"
"os/exec"
"strings"
"sync"
"time"

"github.com/creack/pty"
"github.com/google/uuid"
)

// Shell defines the interface for a shell session
type Shell interface {
Execute(command string) (string, error)
Close() error
}

// ShellSession represents a persistent PTY session
type ShellSession struct {
ID      string
Cmd     *exec.Cmd
Pty     *os.File
outChan chan byte
errChan chan error
mu      sync.Mutex
closed  bool
stop    chan struct{}
}

// NewShellSession spawns a new persistent shell
func NewShellSession(id string) (*ShellSession, error) {
c := exec.Command("bash", "--noprofile", "--norc")

f, err := pty.Start(c)
if err != nil {
return nil, err
}

s := &ShellSession{
ID:      id,
Cmd:     c,
Pty:     f,
outChan: make(chan byte, 8192),
errChan: make(chan error, 1),
stop:    make(chan struct{}),
}

go s.readLoop()

// Disable echo immediately
fmt.Fprintln(f, "stty -echo")

// Initial sync to consume the 'stty -echo' output and any initial prompt
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
_, err = s.executeWithSentinel(ctx, "echo ready")
if err != nil {
s.Close()
return nil, fmt.Errorf("init failed: %v", err)
}

return s, nil
}

func (s *ShellSession) readLoop() {
reader := bufio.NewReader(s.Pty)
for {
b, err := reader.ReadByte()
if err != nil {
select {
case s.errChan <- err:
case <-s.stop:
}
return
}
select {
case s.outChan <- b:
case <-s.stop:
return
}
}
}

func (s *ShellSession) Execute(command string) (string, error) {
s.mu.Lock()
defer s.mu.Unlock()

if s.closed {
return "", fmt.Errorf("session closed")
}

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

return s.executeWithSentinel(ctx, command)
}

func (s *ShellSession) executeWithSentinel(ctx context.Context, command string) (string, error) {
// Drain outChan
for len(s.outChan) > 0 { <-s.outChan }

u := uuid.New().String()
sentinel := "__SENTINEL_" + u + "__"

// Send command and sentinel
_, err := fmt.Fprintln(s.Pty, command + "; echo " + sentinel)
if err != nil {
return "", err
}

var output bytes.Buffer
for {
select {
case b := <-s.outChan:
output.WriteByte(b)
str := output.String()

// Check if sentinel is present
if idx := strings.Index(str, sentinel); idx != -1 {
lastIdx := strings.LastIndex(str, sentinel)
result := str[:lastIdx]
lines := strings.Split(strings.ReplaceAll(result, "\r\n", "\n"), "\n")

// Strip the command echo line if it matches
if len(lines) > 0 && strings.Contains(lines[0], u) {
lines = lines[1:]
}

return strings.TrimSpace(strings.Join(lines, "\n")), nil
}
case err := <-s.errChan:
return output.String(), err
case <-ctx.Done():
return output.String(), ctx.Err()
}
}
}

func (s *ShellSession) Close() error {
s.mu.Lock()
defer s.mu.Unlock()
if s.closed {
return nil
}
s.closed = true
close(s.stop)
s.Pty.Close()
return s.Cmd.Process.Kill()
}

// SessionManager manages multiple shell sessions
type SessionManager struct {
sessions map[string]Shell
mu       sync.RWMutex
Creator  func(id string) (Shell, error)
}

func NewSessionManager() *SessionManager {
return &SessionManager{
sessions: make(map[string]Shell),
Creator: func(id string) (Shell, error) {
return NewShellSession(id)
},
}
}

func (m *SessionManager) GetOrCreate(id string) (Shell, error) {
m.mu.RLock()
s, ok := m.sessions[id]
m.mu.RUnlock()
if ok { return s, nil }

m.mu.Lock()
defer m.mu.Unlock()
if s, ok := m.sessions[id]; ok { return s, nil }

s, err := m.Creator(id)
if err != nil { return nil, err }
m.sessions[id] = s
return s, nil
}

func (m *SessionManager) Cleanup() {
m.mu.Lock()
defer m.mu.Unlock()
for _, s := range m.sessions { s.Close() }
m.sessions = make(map[string]Shell)
}
