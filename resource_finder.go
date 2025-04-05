package konfig

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
)

type root struct {
	path string
	err  error
	once sync.Once
}

var rootInstance = &root{}

// findRootPath searches for a "resources" folder starting from the current working directory
// and moving up the directory tree until found or root is reached.
func findRootPath() (string, error) {
	rootInstance.once.Do(func() {
		dir, err := os.Getwd()
		if err != nil {
			rootInstance.path, rootInstance.err = "", errors.Wrap(err, "failed to get current directory")
			return
		}

		for {
			if stat, err := os.Stat(filepath.Join(dir, "resources")); err == nil && stat.IsDir() {
				rootInstance.path, rootInstance.err = dir, nil
				return
			}

			parentDir := filepath.Dir(dir)
			if parentDir == dir {
				break // Reached the root directory
			}
			dir = parentDir
		}

		rootInstance.path, rootInstance.err = "", errors.New("failed to find root path with 'resources' folder")
	})
	return rootInstance.path, rootInstance.err
}
