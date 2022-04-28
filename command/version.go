package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Ender - A daemon and CLI to securely store your secrets\n")
			fmt.Printf("Version: %s: %s\n", version, architecture)
		},
	}
	version      = "0.0.0"
	architecture = "Linux/amd64"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
