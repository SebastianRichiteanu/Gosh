package main

import (
	"fmt"
	"os"

	"github.com/SebastianRichiteanu/Gosh/internal/autocompleter"
	"github.com/SebastianRichiteanu/Gosh/internal/builtins"
	"github.com/SebastianRichiteanu/Gosh/internal/closer"
	"github.com/SebastianRichiteanu/Gosh/internal/config"
	"github.com/SebastianRichiteanu/Gosh/internal/executor"
	"github.com/SebastianRichiteanu/Gosh/internal/logger"
	"github.com/SebastianRichiteanu/Gosh/internal/prompt"
)

func main() {
	exitCode := run()
	os.Exit(exitCode)
}

func run() int {
	reloadCfgChannel := make(chan bool, 1)

	cfg, err := config.NewConfig(reloadCfgChannel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize config: %v\n", err)
		return 1
	}

	exitChannel := make(chan int, 1)
	builtinCmds := builtins.InitBuiltinCmds(exitChannel, reloadCfgChannel, &cfg.HistoryFile)

	log, err := logger.NewLogger(cfg.LogFile, cfg.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		return 1
	}

	ac := autocompleter.NewAutocompleter(&builtinCmds, cfg, log)
	exec := executor.NewExecutor(&builtinCmds, cfg, log)

	pr, err := prompt.NewPrompt(&builtinCmds, ac, cfg, log)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize prompt: %v\n", err)
		return 1
	}

	c := closer.NewCloser(exitChannel, pr, cfg, log)
	defer c.Recover()
	go c.ListenForSignals()

	return runShellLoop(pr, exec)
}

func runShellLoop(pr *prompt.Prompt, exec *executor.Executor) int {
	var previousInput string

	for {
		cmd, newInput, err := pr.HandlePrompt(previousInput)
		if err != nil {
			panic(err)
		}

		if len(newInput) != 0 {
			previousInput = newInput
			continue
		} else {
			previousInput = ""
		}

		exec.Execute(cmd)
	}
}
