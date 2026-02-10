package cmd

import (
"github.com/spf13/cobra"
)

var (
configPath string
debug      bool
)

var rootCmd = &cobra.Command{
Use:   "hyperagent",
Short: "Hyperagent is an autonomous AI agent and OS companion",
Long:  `Hyperagent is a daemon-based autonomous AI agent that lives in your terminal and helps you manage your system.`,
}

func Execute() error {
return rootCmd.Execute()
}

func init() {
rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "path to config file")
rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")
}
