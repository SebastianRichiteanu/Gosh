package autocompleter

import "strings"

// autoCompletebuiltinCmds finds completions for built-in commands based on the given prefix
func (a *Autocompleter) autoCompletebuiltinCmds(prefix string) []string {
	if a.builtinCmds == nil {
		return nil
	}

	var builtinCmdsSuffixes []string

	for cmd := range *a.builtinCmds {
		after, found := strings.CutPrefix(cmd, prefix)
		if found {
			builtinCmdsSuffixes = append(builtinCmdsSuffixes, after)
		}
	}

	return builtinCmdsSuffixes
}
