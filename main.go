package main

import (
"fmt"
"os/exec"
)

func main() {
fmt.Println("Hyperagent is running...")
}

// ExecuteCommand runs a shell command and returns the output.
func ExecuteCommand(cmd string) (string, error) {
out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
return string(out), err
}
