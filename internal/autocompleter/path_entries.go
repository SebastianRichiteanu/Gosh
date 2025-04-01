package autocompleter

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func (a *Autocompleter) autoCompletePathEntries(prefix string) []string {
	var pathSuffixes []string

	var relPath string
	var err error

	if strings.Contains(prefix, "./") {
		relPath, err = os.Getwd()
		if err != nil {
			a.logger.Warn(fmt.Sprintf("failed to get working dir: %v", err))
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
		afterArr := strings.Split(after, "/")
		after = afterArr[len(afterArr)-1]
		if found {
			pathSuffixes = append(pathSuffixes, after+"/")
		}
	}

	return pathSuffixes
}
