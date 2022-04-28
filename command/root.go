package command

import (
	"context"
	"os"

	"github.com/kaidyth/ender/shared"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:              "ender",
		Short:            "Ender is a place to store your secrets",
		Long:             "A safe and secure, temporary secret storage system.",
		TraverseChildren: true,
	}
)

// Execute runs our root command
func Execute() error {
	shared.NewLogger("TRACE")
	ctx := context.Background()
	return rootCmd.ExecuteContext(ctx)
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
