package autocompleter

import (
	"os"
	"path"
	"strings"
)

func (a *Autocompleter) autoCompletePathEntries(prefix string) []string {
	var pathSuffixes []string

	relPath := ""
	if strings.Contains(prefix, "./") {
		relPath, _ = os.Getwd() // TODO: treat error? maybe once I add debug
		prefix = strings.Trim(prefix, "./")
	} else {
		relPath = path.Dir(prefix)
	}

	files, err := os.ReadDir(relPath)
	if err != nil {
		return pathSuffixes // TODO: treat error? maybe once I add debug
	}

	relevantPrefix := prefix[strings.LastIndex(prefix, "/")+1:]

	for _, file := range files {
		after, found := strings.CutPrefix(file.Name(), relevantPrefix)
		afterArr := strings.Split(after, "/") // TODO: this is sooooo hacky....
		after = afterArr[len(afterArr)-1]
		if found {
			pathSuffixes = append(pathSuffixes, after)
		}
	}

	return pathSuffixes
}
