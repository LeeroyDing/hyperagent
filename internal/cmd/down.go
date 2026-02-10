package cmd

import (
"fmt"
"net/http"
"os"
"path/filepath"
"time"

"github.com/spf13/cobra"
"github.com/LeeroyDing/hyperagent/internal/daemon"
)

var downCmd = &cobra.Command{
Use:   "down",
Short: "Stop the Hyperagent daemon",
Run: func(cmd *cobra.Command, args []string) {
home, _ := os.UserHomeDir()
pidFile := filepath.Join(home, ".hyperagent", "hyperagent.pid")
d := daemon.NewDaemon(pidFile)

pid, err := d.GetPID()
if err != nil {
fmt.Println("Hyperagent daemon is not running.")
return
}

fmt.Printf("Stopping Hyperagent daemon (PID %d)...\n", pid)

// Try graceful shutdown via API first
_, err = http.Post("http://localhost:8080/api/daemon/stop", "application/json", nil)
if err != nil {
fmt.Printf("API shutdown failed, sending SIGTERM to process %d...\n", pid)
proc, err := os.FindProcess(pid)
if err == nil {
proc.Signal(os.Interrupt)
}
}

// Wait for cleanup
for i := 0; i < 5; i++ {
if _, err := os.Stat(pidFile); os.IsNotExist(err) {
fmt.Println("🟢 Hyperagent daemon stopped successfully.")
return
}
time.Sleep(1 * time.Second)
}

fmt.Println("⚠️ Daemon did not stop gracefully. You may need to kill it manually.")
},
}

func init() {
rootCmd.AddCommand(downCmd)
}
