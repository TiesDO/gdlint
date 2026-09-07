package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func FindFileUpwards(fileName string, startingDir string) (string, error) {
	abs, err := filepath.Abs(startingDir)

	if err != nil {
		return "", fmt.Errorf("failed to infer absolute path: %w", err)
	}

	currentDir := abs

	for {
		targetFile := filepath.Join(currentDir, fileName)

		_, err := os.Stat(targetFile)
		if err == nil {
			return targetFile, nil
		}

		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("failed to stat file %s: %w", targetFile, err)
		}

		parentDir := filepath.Dir(currentDir)

		if parentDir == currentDir {
			// if the parent of current is the same as current we have reached
			// the filesystem root
			return "", fmt.Errorf("file %q not found in any parent dir: %w", fileName, os.ErrNotExist)
		}

		currentDir = parentDir
	}
}
