// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"github.com/fsnotify/fsnotify"
	"github.com/osintami/fingerprintz/log"
)

type IFileWatcher interface {
	Add(file string, refresh func()) error
	Listen()
}

type FileWatcher struct {
	watcher *fsnotify.Watcher
	watched map[string]func()
}

func NewFileWatcher() IFileWatcher {
	watcher, _ := fsnotify.NewWatcher()
	return &FileWatcher{
		watcher: watcher,
		watched: make(map[string]func()),
	}
}

func (x *FileWatcher) Add(file string, refresh func()) error {
	err := x.watcher.Add(file)
	if err != nil {
		log.Error().Err(err).Str("component", "watcher").Str("file", file).Msg("add watch")
		return err
	}

	x.watched[file] = refresh
	return nil
}

func (x *FileWatcher) Listen() {
	go func() {
		done := make(chan bool)
		go func() {
			defer close(done)
			for {
				event, ok := <-x.watcher.Events
				if ok && event.Op == fsnotify.Chmod {
					refresh := x.watched[event.Name]
					if refresh != nil {
						refresh()
					}
				}
			}
		}()
		<-done
	}()
}
