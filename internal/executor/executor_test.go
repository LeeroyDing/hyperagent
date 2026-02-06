package executor

import (
"strings"
"testing"
)

func TestShellExecutor_Execute(t *testing.T) {
allowlist := []string{"echo", "ls"}
executor := NewShellExecutor(allowlist)

tests := []struct {
name    string
command string
want    string
wantErr bool
}{
{
name:    "allowed command",
command: "echo hello",
want:    "hello",
wantErr: false,
},
{
name:    "blocked command",
command: "whoami",
want:    "not in the allowlist",
wantErr: true,
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got, err := executor.Execute(tt.command)
if (err != nil) != tt.wantErr {
t.Errorf("ShellExecutor.Execute() error = %v, wantErr %v", err, tt.wantErr)
return
}
if !strings.Contains(strings.ToLower(got), strings.ToLower(tt.want)) && !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tt.want)) {
t.Errorf("ShellExecutor.Execute() got = %v, err = %v, want %v", got, err, tt.want)
}
})
}
}
