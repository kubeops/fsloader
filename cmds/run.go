/*
Copyright AppsCode Inc. and Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmds

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"gomodules.xyz/signals"
	"gomodules.xyz/x/ioutil"
	"k8s.io/klog/v2"
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
			// set up signals, so we handle the first shutdown signal gracefully
			stopCh := signals.SetupSignalHandler()

			var err error
			if meta.PossiblyInCluster() {
				watcher := fsnotify.Watcher{WatchDir: watchDir, Reload: reload}
				err = watcher.Run(stopCh)
			} else {
				watcher := ioutil.Watcher{WatchDir: watchDir, WatchFiles: watchFiles, Reload: reload}
				err = watcher.Run(stopCh)
			}
			if err != nil {
				klog.ErrorS(err, "watcher failed")
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
	klog.Infoln("Output:\n", msg)
	if err != nil {
		klog.Errorln("failed to run cmd")
		return fmt.Errorf("error restarting %v: %v", msg, err)
	}
	return nil
}
