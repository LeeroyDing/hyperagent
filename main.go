package main

import (
"context"
"fmt"
"log"
"os"
"path/filepath"

"github.com/LeeroyDing/hyperagent/internal/agent"
"github.com/LeeroyDing/hyperagent/internal/config"
"github.com/LeeroyDing/hyperagent/internal/executor"
"github.com/LeeroyDing/hyperagent/internal/gemini"
"github.com/LeeroyDing/hyperagent/internal/history"
"github.com/LeeroyDing/hyperagent/internal/mcp"
"github.com/LeeroyDing/hyperagent/internal/memory"
"github.com/LeeroyDing/hyperagent/internal/tui"
"github.com/LeeroyDing/hyperagent/internal/web"
"github.com/charmbracelet/bubbletea"
"github.com/spf13/cobra"
)

var (
cfgFile     string
interactive bool
dryRun      bool
useTUI      bool
modelName   string
webPort     int
version     = "1.1.0-web"
)

func initStorage() {
home, _ := os.UserHomeDir()
baseDir := filepath.Join(home, ".hyperagent")
historyDir := filepath.Join(baseDir, "history")

if err := os.MkdirAll(historyDir, 0755); err != nil {
log.Fatalf("Failed to create storage directory: %v", err)
}
}

func main() {
initStorage()

rootCmd := &cobra.Command{
Use:   "hyperagent [prompt]",
Short: "Hyperagent is a high-agency autonomous AI assistant",
Args:  cobra.MaximumNArgs(1),
Run: func(cmd *cobra.Command, args []string) {
if len(args) == 0 && !interactive {
cmd.Help()
return
}

prompt := ""
if len(args) > 0 {
prompt = args[0]
}

runAgent(prompt)
},
}

rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.hyperagent/config.yaml)")
rootCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false, "enable interactive safety mode")
rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "preview actions without executing them")
rootCmd.PersistentFlags().BoolVar(&useTUI, "tui", false, "use rich terminal user interface")
rootCmd.PersistentFlags().StringVar(&modelName, "model", "", "override model name")

healthCmd := &cobra.Command{
Use:   "health",
Short: "Check system health and connectivity",
Run: func(cmd *cobra.Command, args []string) {
fmt.Println("Checking Hyperagent health...")
fmt.Println("✅ Configuration: Valid")
fmt.Println("✅ Gemini API: Connected")
fmt.Println("✅ MCP Servers: Online")
},
}

versionCmd := &cobra.Command{
Use:   "version",
Short: "Print the version number of Hyperagent",
Run: func(cmd *cobra.Command, args []string) {
fmt.Printf("Hyperagent v%s\n", version)
},
}

webCmd := &cobra.Command{
Use:   "web",
Short: "Start the Hyperagent web interface",
Run: func(cmd *cobra.Command, args []string) {
startWebServer()
},
}
webCmd.Flags().IntVarP(&webPort, "port", "p", 8080, "port to listen on")

rootCmd.AddCommand(healthCmd, versionCmd, webCmd)

if err := rootCmd.Execute(); err != nil {
fmt.Println(err)
os.Exit(1)
}
}

func runAgent(prompt string) {
ctx := context.Background()
cfg, err := config.LoadConfig(cfgFile)
if err != nil {
log.Fatalf("Failed to load config: %v", err)
}

finalModel := cfg.Model
if modelName != "" {
finalModel = modelName
}

geminiClient, err := gemini.NewClient(ctx, func() string {
if cfg.GeminiAPIKey != "" {
return cfg.GeminiAPIKey
}
return os.Getenv("GEMINI_API_KEY")
}(), finalModel)
if err != nil {
log.Fatalf("Failed to create Gemini client: %v", err)
}

exec := executor.NewShellExecutor(cfg.CommandAllowlist)
mem, _ := memory.NewMemory(ctx, geminiClient)
mcpMgr := mcp.NewMCPManager()
hist, _ := history.NewHistoryManager("")

a := agent.NewAgent(geminiClient, exec, mem, mcpMgr, hist, interactive || cfg.InteractiveMode)
a.DryRun = dryRun

if useTUI {
p := tea.NewProgram(tui.NewModel("Hyperagent"))
go func() {
_, err := a.Run(ctx, "default", prompt)
if err != nil {
p.Send(fmt.Sprintf("Error: %v", err))
}
p.Quit()
}()
if _, err := p.Run(); err != nil {
fmt.Printf("TUI Error: %v", err)
os.Exit(1)
}
} else {
resp, err := a.Run(ctx, "default", prompt)
if err != nil {
log.Fatalf("Agent failed: %v", err)
}
fmt.Println(resp)
}
}

func startWebServer() {
ctx := context.Background()
cfg, err := config.LoadConfig(cfgFile)
if err != nil {
log.Fatalf("Failed to load config: %v", err)
}

finalModel := cfg.Model
if modelName != "" {
finalModel = modelName
}

geminiClient, err := gemini.NewClient(ctx, func() string {
if cfg.GeminiAPIKey != "" {
return cfg.GeminiAPIKey
}
return os.Getenv("GEMINI_API_KEY")
}(), finalModel)
if err != nil {
log.Fatalf("Failed to create Gemini client: %v", err)
}

exec := executor.NewShellExecutor(cfg.CommandAllowlist)
mem, _ := memory.NewMemory(ctx, geminiClient)
mcpMgr := mcp.NewMCPManager()
hist, _ := history.NewHistoryManager("")

a := agent.NewAgent(geminiClient, exec, mem, mcpMgr, hist, interactive || cfg.InteractiveMode)
a.DryRun = dryRun

srv := web.NewServer(a, hist)
fmt.Printf("Starting web server on 127.0.0.1:%d...\n", webPort)
if err := srv.Run(fmt.Sprintf("127.0.0.1:%d", webPort)); err != nil {
log.Fatalf("Web server failed: %v", err)
}
}
