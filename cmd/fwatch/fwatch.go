package main

import (
	"flag"
	"strings"
	"time"

	"github.com/vogo/fwatch"
	"github.com/vogo/logger"
)

func main() {
	var (
		dir             = flag.String("dir", "", "directory to watch")
		logLevel        = flag.String("log_level", "", "log level(debug/info)")
		includeSub      = flag.Bool("include_sub", false, "whether include sub-directories")
		fileSuffix      = flag.String("suffix", "", "file suffix to watch")
		inactiveSeconds = flag.Int64("inactive_seconds", 60, "after seconds files is inactive")
	)

	flag.Parse()

	if *dir == "" {
		logger.Fatal("required parameter -dir")
	}

	if strings.EqualFold(*logLevel, "DEBUG") {
		logger.SetLevel(logger.LevelDebug)
	}

	watcher, err := fwatch.NewFileWatcher(*dir, *includeSub, time.Duration(*inactiveSeconds)*time.Second, func(s string) bool {
		return *fileSuffix == "" || strings.HasSuffix(s, *fileSuffix)
	})
	if err != nil {
		logger.Fatal(err)
	}

	defer func() {
		_ = watcher.Stop()
	}()

	go func() {
		for {
			select {
			case <-watcher.Done:
				return
			case f := <-watcher.ActiveChan:
				logger.Infof("--> active file: %s", f.Name)
			case f := <-watcher.InactiveChan:
				logger.Infof("--> inactive file: %s", f.Name)
			case name := <-watcher.RemoveChan:
				logger.Infof("--> remove file: %s", name)
			}
		}
	}()

	if err = watcher.Start(); err != nil {
		logger.Fatal(err)
	}

	select {}
}
