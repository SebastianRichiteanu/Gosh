package autocompleter

import (
	"slices"
	"strings"

	"github.com/SebastianRichiteanu/Gosh/internal/types"
)

type Autocompleter struct {
	builtinCmds *types.CommandMap
}

func NewAutocompleter(builtinCmds *types.CommandMap) *Autocompleter {
	return &Autocompleter{
		builtinCmds: builtinCmds,
	}
}

// Autocomplete generates a list of possible completions for a given prefix
// It combines suggestions from known built-in commands and executable files in the system's PATH
func (a *Autocompleter) Autocomplete(builtinCmds types.CommandMap, prefix string) ([]string, bool) {
	if prefix == "" {
		return nil, false
	}

	if strings.Contains(prefix, " ") || strings.Contains(prefix, "./") || prefix[0] == '/' {
		// If string has space in it probably is an arg of a command
		// TODO: No autocomplete inside quotes
		lastPartPrefix := strings.Split(prefix, " ")

		return a.autoCompletePathEntries(lastPartPrefix[len(lastPartPrefix)-1]), true
	}

	var suffixes []string
	suffixes = append(suffixes, a.autoCompletebuiltinCmds(prefix)...)
	suffixes = append(suffixes, a.autoCompleteExecutables(prefix)...)

	uniqueSuffixes := make(map[string]bool)
	var result []string

	for _, suffix := range suffixes {
		if !uniqueSuffixes[suffix] {
			uniqueSuffixes[suffix] = true
			result = append(result, suffix)
		}
	}

	slices.Sort(result)

	return result, false
}
