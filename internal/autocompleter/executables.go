package autocompleter

import (
	"os"
	"strings"
	"sync"

	"github.com/SebastianRichiteanu/Gosh/internal/types"
)

// autoCompleteExecutables finds completions for executable commands in the system's PATH based on the given prefix
func (a *Autocompleter) autoCompleteExecutables(prefix string) []string {
	path := os.Getenv(types.PathEnvVar)
	directories := strings.Split(path, string(types.PathDelimiter))

	var wg sync.WaitGroup
	suffixesChan := make(chan string)

	// process directories concurrently
	for _, directory := range directories {
		wg.Add(1)
		go func(dir string) {
			defer wg.Done()
			processDirectory(prefix, dir, suffixesChan)
		}(directory)
	}

	go func() {
		wg.Wait()
		close(suffixesChan)
	}()

	var suffixes []string
	// Collect suffixes from channel
	for suffix := range suffixesChan {
		suffixes = append(suffixes, suffix)
	}

	return suffixes
}

// processDirectory searches for file names in a given directory that match the provided prefix
// and sends them to a channel for further processing
func processDirectory(prefix, directory string, suffixesChan chan<- string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if file.Type()&0111 == 0 {
			continue // Skip non-executable files
		}

		after, found := strings.CutPrefix(file.Name(), prefix)
		if found {
			suffixesChan <- after
		}
	}
}
