package command

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"sync"
	"syscall"

	"os/exec"

	"github.com/apex/log"
	api "github.com/kaidyth/ender/client"
	"github.com/kaidyth/ender/server"
	"github.com/spf13/cobra"
)

const SERVER_WAITGROUP_INSTANCES = 1

var (
	wg_count            = 0
	chest               = "default"
	daemonSocketAddress = "default.socket"
	daemonHelperCmd     = &cobra.Command{
		Use: "daemon-helper",
		Run: func(cmd *cobra.Command, args []string) {
			var socket string
			if len(args) == 0 {
				socket = base64.URLEncoding.EncodeToString(api.GenerateRandomBytes(12)) + ".socket"
			} else {
				socket = "default.socket"
			}

			curentDir, _ := os.Getwd()
			ecmd := exec.Command(curentDir+"/ender", "daemon", "--socket="+socket)
			fmt.Printf("export ENDER_SOCKET=%s\n", socket)
			if chest != "" {
				fmt.Printf("export ENDER_CHEST=%s\n", chest)
			}
			ecmd.Start()
			os.Exit(0)
		},
	}

	daemonCmd = &cobra.Command{
		Use: "daemon",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			usr, err := user.Current()
			if err != nil {
				socketAddressPath = "/tmp/"
			} else {
				socketAddressPath = usr.HomeDir + "/.ender/"
			}

			go server.NewServer(ctx, socketAddressPath+daemonSocketAddress)

			var wg sync.WaitGroup
			wg.Add(SERVER_WAITGROUP_INSTANCES)
			wg_count = SERVER_WAITGROUP_INSTANCES
			// Create a signal handler for TERM, INT, and USR1
			var captureSignal = make(chan os.Signal, 1)
			signal.Notify(captureSignal, syscall.SIGINT, syscall.SIGTERM)
			serverSignalHandler(<-captureSignal, &wg)

			// Wait for the goroutines to clearnly exist before ending the server
			wg.Wait()
		},
	}
)

func init() {
	daemonCmd.PersistentFlags().StringVar(&chest, "chest", getenv("ENDER_CHEST", "default"), "The chest name to use.")
	daemonCmd.Flags().StringVar(&daemonSocketAddress, "socket", "default.socket", "The socket name to use")
	rootCmd.AddCommand(daemonCmd)
	rootCmd.AddCommand(daemonHelperCmd)
}

// Signal handler to enable Ender to clean up after itself
func serverSignalHandler(signal os.Signal, wg *sync.WaitGroup) {
	log.Debug("syscall made")
	if signal == syscall.SIGTERM || signal == syscall.SIGINT {
		log.Debug("SIGINT Cleaning up")
		wg.Add(-SERVER_WAITGROUP_INSTANCES)
		// Cleanup the socket
		if err := os.RemoveAll(socketAddressPath + daemonSocketAddress); err != nil {
			log.Debugf("Unable to delete socket: %v", err)
			os.Exit(1)
		}

		// Delete the keyring data
		ring, err := api.GetKeyring(chest)
		if err == nil {
			ring.Remove("key")
		}
		os.Exit(0)
	}
}
