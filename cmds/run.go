package cmds

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync/atomic"

	"github.com/appscode/go/log"
	"github.com/appscode/go/sets"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var (
	fileset    sets.String
	watchFiles []string
	watchDir   string
	reloadCmd  string
)

func NewRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Runs fsloader",
		Run: func(cmd *cobra.Command, args []string) {
			runWatcher()
		},
	}
	cmd.Flags().StringSliceVar(&watchFiles, "watch-files", nil, "Files to be watched")
	cmd.Flags().StringVar(&watchDir, "watch-dir", "", "Dir that contains the files to be watched")
	cmd.Flags().StringVar(&reloadCmd, "reload-cmd", "", "Command to be executed when files are modified")
	return cmd
}

var reloadCount uint64

func incReloadCount() {
	atomic.AddUint64(&reloadCount, 1)
	log.Infoln("reloaded:", atomic.LoadUint64(&reloadCount))
}

func runCmd(path string) error {
	output, err := exec.Command("sh", "-c", path).CombinedOutput()
	msg := fmt.Sprintf("%v", string(output))
	log.Infoln("Output:\n", msg)
	if err != nil {
		log.Errorln("failed to run cmd")
		return fmt.Errorf("error restarting %v: %v", msg, err)
	}
	incReloadCount()
	return nil
}

func runWatcher() {
	fileset = sets.NewString(watchFiles...)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorln(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Infoln("Event: --------------------------------------", event)

				filename := filepath.Clean(event.Name)
				if !fileset.Has(filename) {
					continue
				}

				switch event.Op {
				case fsnotify.Create:
					if err = watcher.Add(filename); err != nil {
						log.Errorln("error:", err)
					}
				case fsnotify.Write:
					if err = printMD5(filename); err == nil {
						log.Errorln("error:", err)
					}
					if err := runCmd(reloadCmd); err != nil {
						log.Errorln(err)
					}
				case fsnotify.Remove, fsnotify.Rename:
					if err = watcher.Remove(filename); err != nil {
						log.Errorln("error:", err)
					}
				}
			case err := <-watcher.Errors:
				log.Errorln("error:", err)
			}
		}
	}()

	for _, filename := range watchFiles {
		if err = watcher.Add(filename); err != nil {
			log.Errorf("error watching file %s. Reason: %s", filename, err)
		}
	}
	if err = watcher.Add(watchDir); err != nil {
		log.Errorf("error watching dir %s. Reason: %s", watchDir, err)
	}
	<-done
}

func printMD5(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	fmt.Printf("%x\n", h.Sum(nil))
	return nil
}
