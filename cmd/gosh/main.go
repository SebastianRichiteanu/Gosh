package main

import (
	"log"
	"os"

	"github.com/SebastianRichiteanu/Gosh/internal/autocompleter"
	"github.com/SebastianRichiteanu/Gosh/internal/builtins"
	"github.com/SebastianRichiteanu/Gosh/internal/executor"
	"github.com/SebastianRichiteanu/Gosh/internal/prompt"
	"github.com/SebastianRichiteanu/Gosh/internal/utils"
)

func main() {
	builtinCmds := builtins.InitBuiltinCmds()

	// TODO: log file path should be in a config :D
	logFile, err := os.OpenFile("gosh.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	autoCompleter := autocompleter.NewAutocompleter(&builtinCmds)
	prompt := prompt.NewPrompt(&builtinCmds, autoCompleter)
	executor := executor.NewExecutor(&builtinCmds)

	var previousInput string

	utils.BlockCtrlC()

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
