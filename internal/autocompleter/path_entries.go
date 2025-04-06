package autocompleter

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/SebastianRichiteanu/Gosh/internal/utils"
)

func (a *Autocompleter) autoCompletePathEntries(prefix string) []string {
	var pathSuffixes []string

	var relPath string
	var err error

	prefix, err = utils.ExpandHomePath(prefix)
	if err != nil {
		a.logger.Warn(fmt.Sprintf("failed to expand home path: %v", err), "prefix", prefix)
	}

	if strings.Contains(prefix, "./") {
		relPath, err = os.Getwd()
		if err != nil {
			a.logger.Warn(fmt.Sprintf("failed to get working dir: %v", err), "prefix", prefix)
		}

		prefix = strings.Trim(prefix, "./")
	} else {
		relPath = path.Dir(prefix)
	}

	files, err := os.ReadDir(relPath)
	if err != nil {
		a.logger.Error(fmt.Sprintf("failed to read dir: %v", err), "path", relPath)
		return pathSuffixes
	}

	relevantPrefix := prefix[strings.LastIndex(prefix, "/")+1:]

	for _, file := range files {
		after, found := strings.CutPrefix(file.Name(), relevantPrefix)
		if found {
			afterArr := strings.Split(after, "/")
			after = afterArr[len(afterArr)-1]

			fullPath := filepath.Join(relPath, file.Name())
			info, err := os.Stat(fullPath)
			if err != nil {
				a.logger.Error(fmt.Sprintf("failed to stat file: %v", err), "path", fullPath)
				continue
			}

			if info.IsDir() {
				pathSuffixes = append(pathSuffixes, after+"/")
			} else {
				pathSuffixes = append(pathSuffixes, after)
			}
		}
	}

	return pathSuffixes
}
