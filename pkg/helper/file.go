package helper

import "os"

func FileExists(file string) bool {
	if file != "" {
		if fileInfo, err := os.Stat(file); err == nil && fileInfo.Mode().IsRegular() {
			return true
		}
	}
	return false
}
