package configmap

import (
	"github.com/fsnotify/fsnotify"
)

// Watcher holds details required to watch for ConfigMap updates
type Watcher struct {
	// Holds the location in which the config file lives
	FilePath string

	// The function to run when an update on the ConfigMap is applied
	OnUpdate func()

	// The function to run when an error occurs
	OnError func(error)

	// The function to run when a fatal error occurs
	// When fatal errors occur, the Watcher is NOT running. Therefore, any ConfigMap updates
	// will be missed.
	OnFatal func(error)
}

// New generates and returns a new Watcher
func New(filePath string, onUpdate func(), onError, onFatal func(error)) *Watcher {
	return &Watcher{
		FilePath: filePath,
		OnUpdate: onUpdate,
		OnError:  onError,
		OnFatal:  onFatal,
	}
}

// Run starts watching for updates to a ConfigMap file
//
// You'll want to call this method in a goroutine or it will block your main thread.
// RunInBackground is also available for just this purpose.
func (w *Watcher) Run() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		w.OnFatal(err)
		return
	}
	defer watcher.Close()

	if err := watcher.Add(w.FilePath); err != nil {
		w.OnFatal(err)
		return
	}

	for {
		select {
		case event := <-watcher.Events:
			// If the ConfigMap is removed during the update, we need to reapply the watcher to
			// refresh the symlink
			if event.Op == fsnotify.Remove {
				watcher.Remove(event.Name)
				watcher.Add(w.FilePath)
			}

			w.OnUpdate()
		case err := <-watcher.Errors:
			w.OnError(err)
		}
	}
}

// RunInBackground starts the Watcher in a new thread
func (w *Watcher) RunInBackground() {
	go w.Run()
}
