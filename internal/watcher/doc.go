// Package watcher provides lightweight file-polling functionality for
// envchain. It watches one or more environment variable files and emits
// an Event on its channel whenever a file's modification time changes.
//
// Usage:
//
//	w := watcher.New(500*time.Millisecond, ".env", ".env.local")
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	go w.Run(ctx)
//
//	for ev := range w.Events() {
//		fmt.Println("changed:", ev.Path)
//		// re-run pipeline ...
//	}
//
// The watcher uses polling rather than OS-level inotify/kqueue so that
// it works uniformly across Linux, macOS, and Windows without additional
// build tags or CGO dependencies.
package watcher
