package autocompleter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SebastianRichiteanu/Gosh/internal/utils"
)

func (a *Autocompleter) autoCompletePathEntries(prefix string) []string {
	var pathSuffixes []string

	expandedPrefix, err := utils.ExpandHomePath(prefix)
	if err != nil {
		a.logger.Warn(fmt.Sprintf("failed to expand home path: %v", err), "prefix", prefix)
		expandedPrefix = prefix
	}

	var (
		dirToRead  string
		basePrefix string
	)

	info, err := os.Stat(expandedPrefix)
	if err == nil && info.IsDir() {
		// If input is a directory and doesn't end with a slash, just return "/"
		if !strings.HasSuffix(prefix, "/") {
			// TODO: this should only happen if the pathsuffixes == 1, eg /mnt/d/Programming/C (C | C++)
			return []string{"/"}
		}

		// if it ends with a slash, list contents
		dirToRead = expandedPrefix
	} else {
		dirToRead = filepath.Dir(expandedPrefix)
		basePrefix = filepath.Base(expandedPrefix)
	}

	files, err := os.ReadDir(dirToRead)
	if err != nil {
		a.logger.Error(fmt.Sprintf("failed to read dir: %v", err), "path", dirToRead)
		return pathSuffixes
	}

	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, basePrefix) {
			suffix := strings.TrimPrefix(name, basePrefix)
			if file.IsDir() {
				suffix += "/"
			}
			pathSuffixes = append(pathSuffixes, suffix)
		}
	}

	return pathSuffixes
}
