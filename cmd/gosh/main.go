package main

import (
	"github.com/SebastianRichiteanu/Gosh/internal/builtins"
	"github.com/SebastianRichiteanu/Gosh/internal/executor"
	"github.com/SebastianRichiteanu/Gosh/internal/prompt"
)

func main() {
	knownCmds := builtins.InitBuiltins()

	var previousInput string // TODO: scrap this and use history
	var history []string

	for {
		prompt, newInput, newHistory, err := prompt.Prompt(knownCmds, previousInput, history)
		if err != nil {
			panic(err)
		}

		if len(newInput) != 0 {
			previousInput = newInput
			continue
		} else {
			previousInput = ""
		}

		history = newHistory

		executor.Exec(prompt, knownCmds)
	}
}
