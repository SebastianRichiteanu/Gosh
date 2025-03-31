package main

import (
	"fmt"
	"os"

	"github.com/SebastianRichiteanu/Gosh/internal/autocompleter"
	"github.com/SebastianRichiteanu/Gosh/internal/builtins"
	"github.com/SebastianRichiteanu/Gosh/internal/config"
	"github.com/SebastianRichiteanu/Gosh/internal/executor"
	"github.com/SebastianRichiteanu/Gosh/internal/logger"
	"github.com/SebastianRichiteanu/Gosh/internal/prompt"
	"github.com/SebastianRichiteanu/Gosh/internal/utils"
)

func main() {
	builtinCmds, reloadCfgChannel := builtins.InitBuiltinCmds()

	cfg := config.NewConfig(reloadCfgChannel)

	logger, err := logger.NewLogger(cfg.LogFile, cfg.LogLevel) // TODO: this is not dynamic
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	autoCompleter := autocompleter.NewAutocompleter(&builtinCmds, cfg, logger)
	prompt := prompt.NewPrompt(&builtinCmds, autoCompleter, cfg, logger)
	executor := executor.NewExecutor(&builtinCmds, cfg, logger)

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
