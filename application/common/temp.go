package common

import (
	"fmt"
	"os"
)

// GetTempDir returns temporary directory that can be used by install / update scenarios
func GetTempDir() (dir string, cleanup func(), err error) {
	dir, err = os.MkdirTemp("", "myapps")
	if err != nil {
		err = fmt.Errorf("failed to create temp directory for this scenario: %w", err)
	}
	cleanup = func() {
		os.RemoveAll(dir)
	}
	return
}
