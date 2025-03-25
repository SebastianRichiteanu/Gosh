package main

import (
	"github.com/SebastianRichiteanu/Gosh/internal/builtins"
	"github.com/SebastianRichiteanu/Gosh/internal/executor"
	"github.com/SebastianRichiteanu/Gosh/internal/prompt"
)

func main() {
	knownCmds := builtins.InitBuiltins()

	var previousInput string

	for {
		prompt, newInput, err := prompt.Prompt(knownCmds, previousInput)
		if err != nil {
			panic(err)
		}

		if len(newInput) != 0 {
			previousInput = newInput
			continue
		} else {
			previousInput = ""
		}

		executor.Exec(prompt, knownCmds)
	}
}
