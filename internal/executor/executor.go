package executor

import (
"fmt"
"log/slog"
"strings"
)

type Executor interface {
	Execute(sessionID, command string) (string, error)
}

type ShellExecutor struct {
Allowlist []string
Manager   *SessionManager
}

func NewShellExecutor(allowlist []string) *ShellExecutor {
return &ShellExecutor{
Allowlist: allowlist,
Manager:   NewSessionManager(),
}
}

func (e *ShellExecutor) Execute(sessionID, command string) (string, error) {
cmdParts := strings.Fields(command)
if len(cmdParts) == 0 {
return "", fmt.Errorf("empty command")
}

baseCmd := cmdParts[0]
if len(e.Allowlist) > 0 {
allowed := false
for _, a := range e.Allowlist {
if baseCmd == a {
allowed = true
break
}
}
if !allowed {
slog.Warn("Command blocked by allowlist", "command", baseCmd)
return "", fmt.Errorf("command '%s' is not in the allowlist", baseCmd)
}
}

slog.Debug("Executing shell command in session", "session", sessionID, "command", command)

session, err := e.Manager.GetOrCreate(sessionID)
if err != nil {
return "", fmt.Errorf("failed to get shell session: %v", err)
}

return session.Execute(command)
}

func (e *ShellExecutor) Cleanup() {
e.Manager.Cleanup()
}
