package daemon

import (
"fmt"
"os"
"path/filepath"
"strconv"
"syscall"
)

// Daemon handles PID file management and process lifecycle
type Daemon struct {
PIDFile string
}

func NewDaemon(pidFile string) *Daemon {
return &Daemon{PIDFile: pidFile}
}

// Lock creates a PID file or returns an error if another instance is running
func (d *Daemon) Lock() error {
// Ensure directory exists
dir := filepath.Dir(d.PIDFile)
if err := os.MkdirAll(dir, 0755); err != nil {
return fmt.Errorf("failed to create daemon directory: %w", err)
}

// Check if PID file exists
if _, err := os.Stat(d.PIDFile); err == nil {
data, err := os.ReadFile(d.PIDFile)
if err == nil {
pid, _ := strconv.Atoi(string(data))
process, err := os.FindProcess(pid)
if err == nil {
// Check if process is actually running
err := process.Signal(syscall.Signal(0))
if err == nil {
return fmt.Errorf("daemon is already running with PID %d", pid)
}
}
}
}

// Write current PID to file
pid := os.Getpid()
if err := os.WriteFile(d.PIDFile, []byte(strconv.Itoa(pid)), 0644); err != nil {
return fmt.Errorf("failed to write PID file: %w", err)
}

return nil
}

// Unlock removes the PID file
func (d *Daemon) Unlock() error {
if err := os.Remove(d.PIDFile); err != nil && !os.IsNotExist(err) {
return fmt.Errorf("failed to remove PID file: %w", err)
}
return nil
}

// GetPID returns the PID from the PID file
func (d *Daemon) GetPID() (int, error) {
data, err := os.ReadFile(d.PIDFile)
if err != nil {
return 0, err
}
return strconv.Atoi(string(data))
}
