package store

import "os"

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
