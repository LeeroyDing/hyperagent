package cmd

import (
"fmt"
"net/http"
"os"
"path/filepath"

"github.com/spf13/cobra"
"github.com/LeeroyDing/hyperagent/internal/daemon"
)

var statusCmd = &cobra.Command{
Use:   "status",
Short: "Check the status of the Hyperagent daemon",
Run: func(cmd *cobra.Command, args []string) {
home, _ := os.UserHomeDir()
pidFile := filepath.Join(home, ".hyperagent", "hyperagent.pid")
d := daemon.NewDaemon(pidFile)

pid, err := d.GetPID()
if err != nil {
fmt.Println("🔴 Hyperagent daemon is not running (no PID file found).")
return
}

// Try to ping the API
resp, err := http.Get("http://localhost:8080/api/daemon/status")
if err != nil {
fmt.Printf("🟡 Hyperagent daemon (PID %d) is running but API is unreachable: %v\n", pid, err)
return
}
defer resp.Body.Close()

if resp.StatusCode == http.StatusOK {
fmt.Printf("🟢 Hyperagent daemon (PID %d) is running and healthy.\n", pid)
} else {
fmt.Printf("🟡 Hyperagent daemon (PID %d) returned status %d.\n", pid, resp.StatusCode)
}
},
}

func init() {
rootCmd.AddCommand(statusCmd)
}
