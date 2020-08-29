# configmap

A simple watcher for handling updates to your Kubernetes ConfigMap files.

## Installing

```
go get github.com/codykaup/configmap
```

## Usage

The goal is to build a `Watcher` struct with simple names to easily determine what will happen when the ConfigMap updates or gets an error.

```go
import (
	"fmt"

	"github.com/codykaup/configmap"
)

func main() {
	w := &configmap.Watcher{
		FilePath: "path/to/configmap",
		OnUpdate: func() {
			fmt.Println("configmap updated")
		},
		OnError: func(err error) {
			fmt.Println("error found")
		},
		OnFatal: func(err error) {
			panic("fatal error found")
		},
	}

	go w.Run()
	// or
	w.RunInBackground()
}
```

Keep in mind, `OnFatal` means the `Watcher` failed to get started and, thus, the ConfigMap is not being watched. You'll likely want to handle this error or fatal out as well.
