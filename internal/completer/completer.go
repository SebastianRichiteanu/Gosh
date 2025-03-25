package completer

import (
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/SebastianRichiteanu/Gosh/internal/types"
)

// FindLongestPrefix takes a list of command suffixes and returns the longest common prefix among them
func FindLongestPrefix(cmds []string) string {
	common := ""
	isCommon := true

	for i := 0; isCommon; i++ {
		var current byte
		for j := 0; j < len(cmds); j++ {
			if i >= len(cmds[j]) {
				isCommon = false
				break
			} else if j == 0 {
				current = cmds[j][i]
			} else if current != cmds[j][i] {
				isCommon = false
				break
			}
		}
		if isCommon {
			common += string(current)
		}
	}
	return common
}

// Autocomplete generates a list of possible completions for a given prefix
// It combines suggestions from known built-in commands and executable files in the system's PATH
func Autocomplete(knownCmds types.CommandMap, prefix string) []string {
	var suffixes []string

	if prefix == "" {
		return suffixes
	}

	suffixes = append(suffixes, autoCompleteKnownCmds(knownCmds, prefix)...)
	suffixes = append(suffixes, autoCompleteExecutables(prefix)...)

	uniqueSuffixes := make(map[string]bool)
	var result []string

	for _, suffix := range suffixes {
		if !uniqueSuffixes[suffix] {
			uniqueSuffixes[suffix] = true
			result = append(result, suffix)
		}
	}

	slices.Sort(result)

	return result
}

// autoCompleteKnownCmds finds completions for built-in commands based on the given prefix
func autoCompleteKnownCmds(knownCmds types.CommandMap, prefix string) []string {
	var knownCmdsSuffixes []string

	for cmd := range knownCmds {
		after, found := strings.CutPrefix(cmd, prefix)
		if found {
			knownCmdsSuffixes = append(knownCmdsSuffixes, after)
		}
	}

	return knownCmdsSuffixes
}

// autoCompleteExecutables finds completions for executable commands in the system's PATH based on the given prefix
func autoCompleteExecutables(prefix string) []string {
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

		// if file.Type()&0111 == 0 {
		// 	continue // Skip non-executable files
		// }
		// TODO:: this is really slow :(

		after, found := strings.CutPrefix(file.Name(), prefix)
		if found {
			suffixesChan <- after
		}
	}
}
