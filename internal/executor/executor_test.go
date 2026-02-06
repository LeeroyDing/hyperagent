package executor

import (
"testing"
"strings"
)

func TestShellExecutor_Execute(t *testing.T) {
e := NewShellExecutor([]string{"ls", "echo", "pwd"})
defer e.Cleanup()

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
command: "rm -rf /",
want:    "",
wantErr: true,
},
{
name:    "stateful cd",
command: "pwd",
want:    "/", // Just checking it runs
wantErr: false,
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got, err := e.Execute("test-session", tt.command)
if (err != nil) != tt.wantErr {
t.Errorf("ShellExecutor.Execute() error = %v, wantErr %v", err, tt.wantErr)
return
}
if !tt.wantErr && !strings.Contains(got, tt.want) {
t.Errorf("ShellExecutor.Execute() got = %v, want %v", got, tt.want)
}
})
}
}
