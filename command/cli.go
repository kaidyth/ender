package command

import (
	"fmt"
	"os"
	"os/user"

	api "github.com/kaidyth/ender/client"
	"github.com/spf13/cobra"
)

var (
	socketAddressPath string
	socketAddress     string
	cliCmd            = &cobra.Command{
		Use:              "cli",
		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	getCmd = &cobra.Command{
		Use:        "get [flags] key",
		Args:       cobra.MinimumNArgs(1),
		ArgAliases: []string{"key"},
		PreRun:     CliCommandPreRun,
		Run: func(cmd *cobra.Command, args []string) {
			result, err := api.Get(socketAddressPath+socketAddress, chest, args[0])
			if err == nil {
				fmt.Printf("%s", result)
			} else {
				os.Exit(1)
			}
		},
	}

	setCmd = &cobra.Command{
		Use:        "set [flags] key value",
		Args:       cobra.MinimumNArgs(2),
		ArgAliases: []string{"key", "value"},
		PreRun:     CliCommandPreRun,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := api.Set(socketAddressPath+socketAddress, chest, args[0], args[1])
			if err == nil {
				os.Exit(0)
			}
			os.Exit(1)
		},
	}

	deleteCmd = &cobra.Command{
		Use:        "delete key",
		Args:       cobra.MinimumNArgs(1),
		ArgAliases: []string{"key"},
		PreRun:     CliCommandPreRun,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := api.Del(socketAddressPath+socketAddress, chest, args[0])
			if err == nil {
				os.Exit(0)
			}
			os.Exit(1)
		},
	}

	existsCmd = &cobra.Command{
		Use:        "exists",
		Args:       cobra.MinimumNArgs(1),
		ArgAliases: []string{"key"},
		PreRun:     CliCommandPreRun,
		Run: func(cmd *cobra.Command, args []string) {
			result, err := api.Exists(socketAddressPath+socketAddress, chest, args[0])
			if err == nil && result {
				os.Exit(0)
			}
			os.Exit(1)
		},
	}
)

func CliCommandPreRun(cmd *cobra.Command, args []string) {
	usr, err := user.Current()
	if err != nil {
		socketAddressPath = "/tmp/"
	} else {
		socketAddressPath = usr.HomeDir + "/.ender/"
	}
}

func init() {
	cliCmd.PersistentFlags().StringVar(&chest, "chest", getenv("ENDER_CHEST", "default"), "The chest name to use.")
	cliCmd.PersistentFlags().StringVar(&socketAddress, "socket", getenv("ENDER_SOCKET", "default.socket"), "The socket name to use")

	rootCmd.AddCommand(cliCmd)
	cliCmd.AddCommand(getCmd)
	cliCmd.AddCommand(setCmd)
	cliCmd.AddCommand(deleteCmd)
	cliCmd.AddCommand(existsCmd)
}
