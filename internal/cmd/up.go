package cmd

import (
"context"
"log/slog"
"os"
"os/exec"
"os/signal"
"path/filepath"
"syscall"

"github.com/spf13/cobra"
"github.com/LeeroyDing/hyperagent/internal/agent"
"github.com/LeeroyDing/hyperagent/internal/config"
"github.com/LeeroyDing/hyperagent/internal/daemon"
"github.com/LeeroyDing/hyperagent/internal/executor"
"github.com/LeeroyDing/hyperagent/internal/gemini"
"github.com/LeeroyDing/hyperagent/internal/history"
"github.com/LeeroyDing/hyperagent/internal/memory"
"github.com/LeeroyDing/hyperagent/internal/mcp"
"github.com/LeeroyDing/hyperagent/internal/web"
)

var daemonize bool

var upCmd = &cobra.Command{
Use:   "up",
Short: "Start the Hyperagent daemon",
Run: func(cmd *cobra.Command, args []string) {
home, _ := os.UserHomeDir()
workDir := filepath.Join(home, ".hyperagent")
pidFile := filepath.Join(workDir, "hyperagent.pid")
logFile := filepath.Join(workDir, "hyperagent.log")

if daemonize {
// Ensure workdir exists
os.MkdirAll(workDir, 0755)

// Prepare command to run in background
newArgs := []string{}
for _, arg := range os.Args[1:] {
if arg != "--daemon" && arg != "-d" {
newArgs = append(newArgs, arg)
}
}

cmd := exec.Command(os.Args[0], newArgs...)

// Redirect output to log file
f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
if err != nil {
slog.Error("Failed to open log file", "error", err)
os.Exit(1)
}
cmd.Stdout = f
cmd.Stderr = f

if err := cmd.Start(); err != nil {
slog.Error("Failed to start daemon", "error", err)
os.Exit(1)
}

slog.Info("Hyperagent daemon started in background", "pid", cmd.Process.Pid, "log", logFile)
os.Exit(0)
}

d := daemon.NewDaemon(pidFile)
if err := d.Lock(); err != nil {
slog.Error("Failed to lock PID file", "error", err)
os.Exit(1)
}
defer d.Unlock()

// Setup logging
level := slog.LevelInfo
if debug {
level = slog.LevelDebug
}
slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})))

cfg, err := config.LoadConfig(configPath)
if err != nil {
if os.IsNotExist(err) {
cfg, err = config.RunOOBE()
if err != nil {
slog.Error("Failed to run setup", "error", err)
os.Exit(1)
}
} else {
slog.Error("Failed to load config", "error", err)
os.Exit(1)
}
}

ctx := context.Background()
gClient, err := gemini.NewClient(ctx, cfg.GeminiAPIKey, cfg.Model)
if err != nil {
slog.Error("Failed to initialize Gemini client", "error", err)
os.Exit(1)
}

executor := executor.NewShellExecutor(cfg.CommandAllowlist)
mem, err := memory.NewMemory(ctx, gClient, "")
if err != nil {
slog.Error("Failed to initialize memory", "error", err)
os.Exit(1)
}

mcpMgr := mcp.NewMCPManager()
historyMgr, err := history.NewHistoryManager(history.GetDefaultHistoryDir())
if err != nil {
slog.Error("Failed to initialize history manager", "error", err)
os.Exit(1)
}

a := agent.NewAgent(gClient, executor, mem, mcpMgr, historyMgr, cfg.InteractiveMode)

srv := web.NewServer(a, historyMgr, mem, d)

// Handle cleanup on exit
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)
go func() {
<-c
slog.Info("Shutting down...")
executor.Cleanup()
d.Unlock()
os.Exit(0)
}()

addr := "127.0.0.1:8080"
slog.Info("Starting Hyperagent daemon", "addr", addr, "pid", os.Getpid())
if err := srv.Run(addr); err != nil {
slog.Error("Daemon API error", "error", err)
os.Exit(1)
}
},
}

func init() {
upCmd.Flags().BoolVarP(&daemonize, "daemon", "d", false, "Run in background as a daemon")
rootCmd.AddCommand(upCmd)
}
