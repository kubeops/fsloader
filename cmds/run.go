package cmds

import (
	"fmt"
	"os/exec"

	"github.com/appscode/go/ioutil"
	"github.com/appscode/go/log"
	"github.com/appscode/go/signals"
	"github.com/spf13/cobra"
	"kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/tools/fsnotify"
)

var (
	watchFiles []string
	watchDir   string
	reloadCmd  string
)

func NewRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Runs fsloader",
		Run: func(cmd *cobra.Command, args []string) {
			// set up signals so we handle the first shutdown signal gracefully
			stopCh := signals.SetupSignalHandler()

			if meta.PossiblyInCluster() {
				watcher := fsnotify.Watcher{WatchDir: watchDir, Reload: reload}
				watcher.Run(stopCh)
			} else {
				watcher := ioutil.Watcher{WatchDir: watchDir, WatchFiles: watchFiles, Reload: reload}
				watcher.Run(stopCh)
			}

			<-stopCh
		},
	}
	cmd.Flags().StringSliceVar(&watchFiles, "watch-files", nil, "Files to be watched")
	cmd.Flags().StringVar(&watchDir, "watch-dir", "", "Dir that contains the files to be watched")
	cmd.Flags().StringVar(&reloadCmd, "reload-cmd", "", "Command to be executed when files are modified")
	return cmd
}

func reload() error {
	output, err := exec.Command("sh", "-c", reloadCmd).CombinedOutput()
	msg := fmt.Sprintf("%v", string(output))
	log.Infoln("Output:\n", msg)
	if err != nil {
		log.Errorln("failed to run cmd")
		return fmt.Errorf("error restarting %v: %v", msg, err)
	}
	return nil
}
