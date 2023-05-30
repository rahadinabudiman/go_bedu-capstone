package helpers

import (
	"path/filepath"
	"strings"
)

func GetFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.ToLower(ext)
}
