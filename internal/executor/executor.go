package executor

import (
"fmt"
"log/slog"
"os/exec"
"strings"
)

type ShellExecutor struct {
Allowlist []string
}

func NewShellExecutor(allowlist []string) *ShellExecutor {
return &ShellExecutor{Allowlist: allowlist}
}

func (e *ShellExecutor) Execute(command string) (string, error) {
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

slog.Debug("Executing shell command", "command", command)
cmd := exec.Command("sh", "-c", command)
out, err := cmd.CombinedOutput()
return strings.TrimSpace(string(out)), err
}
