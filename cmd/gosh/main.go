package main

import (
	"fmt"
	"os"

	"github.com/SebastianRichiteanu/Gosh/internal/autocompleter"
	"github.com/SebastianRichiteanu/Gosh/internal/builtins"
	"github.com/SebastianRichiteanu/Gosh/internal/executor"
	"github.com/SebastianRichiteanu/Gosh/internal/logger"
	"github.com/SebastianRichiteanu/Gosh/internal/prompt"
	"github.com/SebastianRichiteanu/Gosh/internal/utils"
)

func main() {
	// TODO: create config

	builtinCmds := builtins.InitBuiltinCmds()

	logger, err := logger.NewLogger("gosh.log", "INFO")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	autoCompleter := autocompleter.NewAutocompleter(&builtinCmds, logger)
	prompt := prompt.NewPrompt(&builtinCmds, autoCompleter, logger)
	executor := executor.NewExecutor(&builtinCmds, logger)

	utils.BlockCtrlC()

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
