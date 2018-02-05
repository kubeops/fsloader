package cmds

import (
	"fmt"
	"os/exec"

	"github.com/appscode/go/ioutil"
	"github.com/appscode/go/log"
	"github.com/appscode/go/signals"
	"github.com/spf13/cobra"
)

func NewRunCmd() *cobra.Command {
	var reloadCmd string
	watcher := ioutil.Watcher{
		Reload: func() error {
			output, err := exec.Command("sh", "-c", reloadCmd).CombinedOutput()
			msg := fmt.Sprintf("%v", string(output))
			log.Infoln("Output:\n", msg)
			if err != nil {
				log.Errorln("failed to run cmd")
				return fmt.Errorf("error restarting %v: %v", msg, err)
			}
			return nil
		},
	}

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Runs fsloader",
		Run: func(cmd *cobra.Command, args []string) {
			// set up signals so we handle the first shutdown signal gracefully
			stopCh := signals.SetupSignalHandler()

			watcher.Run(stopCh)

			<-stopCh
		},
	}
	cmd.Flags().StringSliceVar(&watcher.WatchFiles, "watch-files", nil, "Files to be watched")
	cmd.Flags().StringVar(&watcher.WatchDir, "watch-dir", "", "Dir that contains the files to be watched")
	cmd.Flags().StringVar(&reloadCmd, "reload-cmd", "", "Command to be executed when files are modified")
	return cmd
}
