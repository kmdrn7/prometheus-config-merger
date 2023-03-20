package watcher

import (
	"log"
)

func (w *Watcher) DryRun() error {
	log.Println("running in dryrun mode")
	w.setupConfig()
	w.setupWatchedResources()
	return nil
}
