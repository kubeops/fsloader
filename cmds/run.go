package cmds

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sync/atomic"

	"github.com/appscode/go/log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var (
	mountFile, mountDir, bashFile string
)

func NewRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Runs fsloader",
		Run: func(cmd *cobra.Command, args []string) {
			fileWatchTest()
		},
	}
	cmd.Flags().StringVar(&mountFile, "mount-file", "", "Volume location where the file will be mounted")
	cmd.Flags().StringVar(&mountDir, "mount-dir", "", "Volume location where the file will be mounted")
	cmd.Flags().StringVar(&bashFile, "boot-cmd", "", "Bash script that will be run on every change of the file")
	return cmd
}

var updateReceived, mountPerformed uint64

func incUpdateReceivedCounter() {
	atomic.AddUint64(&updateReceived, 1)
	log.Infoln("Update Received:", atomic.LoadUint64(&updateReceived))
}

func incMountCounter() {
	atomic.AddUint64(&mountPerformed, 1)
	log.Infoln("Mount Performed:", atomic.LoadUint64(&mountPerformed))
}

func runCmd(path string) error {
	log.Infoln("calling boot file to execute")
	output, err := exec.Command("sh", "-c", path).CombinedOutput()
	msg := fmt.Sprintf("%v", string(output))
	log.Infoln("Output:\n", msg)
	if err != nil {
		log.Errorln("failed to run cmd")
		return fmt.Errorf("error restarting %v: %v", msg, err)
	}
	log.Infoln("boot file executed")
	return nil
}

func fileWatchTest() {
	fmt.Println("-----------------File Watch Test Began----------------")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error Occured ***************")
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Infoln("Event:-------------------------------------------------------", event, reflect.TypeOf(event), event.String())
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Infoln("modified file:", event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Infoln("Removed file:---------------", event.Name)

					if filepath.Clean(event.Name) == mountFile {
						err = printMD5(mountFile)
						if err != nil {
							log.Errorln("fffffffffffffffffff Error", err)
						}
						err = watcher.Add(mountFile)
						if err != nil {
							log.Errorln("wwwwwwwwwwwwwwwwwww Error", err)
						}
					}
				}
			case err := <-watcher.Errors:
				log.Infoln("error:", err)
			}
		}
	}()

	err = watcher.Add(mountFile)
	if err != nil {
		log.Fatalln("1st Error", err)
	}
	err = watcher.Add(mountDir)
	if err != nil {
		log.Fatalln("2nd Error", err)
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
