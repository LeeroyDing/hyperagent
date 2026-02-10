package cmd

import (
"fmt"
"github.com/spf13/cobra"
)

var version = "v0.0.13"

var versionCmd = &cobra.Command{
Use:   "version",
Short: "Print the version number of Hyperagent",
Run: func(cmd *cobra.Command, args []string) {
fmt.Printf("Hyperagent %s\n", version)
},
}

func init() {
rootCmd.AddCommand(versionCmd)
}
