package main

import (
"context"
"flag"
"fmt"
"log/slog"
"os"
"os/signal"
"syscall"

"github.com/LeeroyDing/hyperagent/internal/agent"
"github.com/LeeroyDing/hyperagent/internal/config"
"github.com/LeeroyDing/hyperagent/internal/executor"
"github.com/LeeroyDing/hyperagent/internal/gemini"
"github.com/LeeroyDing/hyperagent/internal/history"
"github.com/LeeroyDing/hyperagent/internal/memory"
"github.com/LeeroyDing/hyperagent/internal/mcp"
"github.com/LeeroyDing/hyperagent/internal/web"
)

var version = "v0.0.12"

func main() {
configPath := flag.String("config", "", "Path to config file")
interactive := flag.Bool("interactive", false, "Enable interactive mode")
debug := flag.Bool("debug", false, "Enable debug logging")
webMode := flag.Bool("web", false, "Start web server")
port := flag.Int("port", 3001, "Web server port")
showVersion := flag.Bool("version", false, "Show version")
flag.Parse()

if *showVersion {
fmt.Printf("Hyperagent %s\n", version)
return
}

level := slog.LevelInfo
if *debug {
level = slog.LevelDebug
}
slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})))

cfg, err := config.LoadConfig(*configPath)
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

exec := executor.NewShellExecutor(cfg.CommandAllowlist)

mem, err := memory.NewMemory(ctx, gClient)
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

a := agent.NewAgent(gClient, exec, mem, mcpMgr, historyMgr, *interactive || cfg.InteractiveMode)

// Handle cleanup on exit
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)
go func() {
<-c
slog.Info("Shutting down...")
exec.Cleanup()
os.Exit(0)
}()

if *webMode {
srv := web.NewServer(a, historyMgr, mem)
addr := fmt.Sprintf("127.0.0.1:%d", *port)
slog.Info("Starting web server", "addr", addr)
if err := srv.Run(addr); err != nil {
slog.Error("Web server error", "error", err)
os.Exit(1)
}
} else if len(flag.Args()) > 0 {
prompt := flag.Args()[0]
resp, err := a.Run(ctx, "default", prompt)
if err != nil {
slog.Error("Agent error", "error", err)
os.Exit(1)
}
fmt.Println(resp)
} else {
fmt.Printf("Hyperagent %s - Ready\n", version)
}
}
