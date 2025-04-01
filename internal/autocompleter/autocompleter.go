package autocompleter

import (
	"slices"

	"github.com/SebastianRichiteanu/Gosh/internal/config"
	"github.com/SebastianRichiteanu/Gosh/internal/logger"
	"github.com/SebastianRichiteanu/Gosh/internal/types"
)

type Autocompleter struct {
	cfg         *config.Config
	builtinCmds *types.CommandMap
	logger      *logger.Logger
}

func NewAutocompleter(builtinCmds *types.CommandMap, cfg *config.Config, logger *logger.Logger) *Autocompleter {
	return &Autocompleter{
		cfg:         cfg,
		builtinCmds: builtinCmds,
		logger:      logger,
	}
}

// Autocomplete generates a list of possible completions for a given prefix
// It combines suggestions from known built-in commands and executable files in the system's PATH
func (a *Autocompleter) Autocomplete(builtinCmds types.CommandMap, input string) []string {
	if !a.cfg.EnableAutoComplete {
		return nil
	}

	if input == "" {
		return nil
	}

	pathSuffixes := a.autoCompletePathEntries(input)
	if len(pathSuffixes) != 0 {
		return pathSuffixes
	}

	var suffixes []string
	suffixes = append(suffixes, a.autoCompletebuiltinCmds(input)...)
	suffixes = append(suffixes, a.autoCompleteExecutables(input)...)

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
