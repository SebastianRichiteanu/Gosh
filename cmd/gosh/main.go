package main

import (
	"github.com/SebastianRichiteanu/Gosh/internal/builtins"
	"github.com/SebastianRichiteanu/Gosh/internal/completer"
	"github.com/SebastianRichiteanu/Gosh/internal/executor"
	"github.com/SebastianRichiteanu/Gosh/internal/prompt"
)

func main() {
	knownCmds := builtins.InitBuiltins()

	autoCompleter := completer.NewAutocompleter(&knownCmds)
	prompt := prompt.NewPrompt(&knownCmds, autoCompleter)
	executor := executor.NewExecutor(&knownCmds)

	var previousInput string

	for {
		prompt, newInput, err := prompt.HandlePrompt(previousInput)
		if err != nil {
			panic(err)
		}

		if len(newInput) != 0 {
			previousInput = newInput
			continue
		} else {
			previousInput = ""
		}

		executor.Execute(prompt)
	}
}
