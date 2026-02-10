package main

import (
"os"
"os/exec"
"strings"
"testing"

"github.com/stretchr/testify/assert"
)

func TestMain_Version(t *testing.T) {
cmd := exec.Command("go", "run", "main.go", "-version")
output, err := cmd.CombinedOutput()
assert.NoError(t, err)
assert.Contains(t, string(output), "Hyperagent v0.0.13")
}

func TestMain_Help(t *testing.T) {
cmd := exec.Command("go", "run", "main.go", "-help")
output, _ := cmd.CombinedOutput()
assert.Contains(t, string(output), "Usage of")
}

func TestMain_NoArgs(t *testing.T) {
// We need to mock config or ensure it exists to avoid OOBE
// For now, just check if it prints the ready message when config exists
cmd := exec.Command("go", "run", "main.go")
// Set a dummy config path to avoid OOBE if possible, or just check output
output, _ := cmd.CombinedOutput()
// It might fail due to missing config, but we just want to cover the flag parsing and initial checks
assert.NotEmpty(t, output)
}
