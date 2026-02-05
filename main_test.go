package main

import (
"strings"
"testing"
)

func TestExecuteCommand(t *testing.T) {
tests := []struct {
name    string
command string
want    string
wantErr bool
}{
{"echo test", "echo hello", "hello", false},
{"invalid command", "nonexistentcommand", "", true},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got, err := ExecuteCommand(tt.command)
if (err != nil) != tt.wantErr {
t.Errorf("ExecuteCommand() error = %v, wantErr %v", err, tt.wantErr)
return
}
if !tt.wantErr && !strings.Contains(got, tt.want) {
t.Errorf("ExecuteCommand() got = %v, want %v", got, tt.want)
}
})
}
}
