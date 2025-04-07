package prompt

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/SebastianRichiteanu/Gosh/internal/autocompleter"
	"github.com/SebastianRichiteanu/Gosh/internal/config"
	"github.com/SebastianRichiteanu/Gosh/internal/logger"
	"github.com/SebastianRichiteanu/Gosh/internal/types"
	"github.com/mattn/go-tty"
)

var (
	errUnterminatedSingleQuote = errors.New("unterminated single quotes")
	errUnterminatedDoubleQuote = errors.New("unterminated double quotes")
)

type Prompt struct {
	cfg           *config.Config
	builtinCmds   *types.CommandMap
	autocompleter *autocompleter.Autocompleter
	logger        *logger.Logger

	tty *tty.TTY

	osSignalsChan chan os.Signal
	runeChan      chan rune
	errChan       chan error

	history      []string
	historyIndex int
}

func NewPrompt(builtinCmds *types.CommandMap, autocompleter *autocompleter.Autocompleter, cfg *config.Config, logger *logger.Logger) (*Prompt, error) {
	p := Prompt{
		cfg:           cfg,
		builtinCmds:   builtinCmds,
		autocompleter: autocompleter,
		logger:        logger,

		osSignalsChan: make(chan os.Signal, 1),
		runeChan:      make(chan rune),
		errChan:       make(chan error),

		history:      []string{},
		historyIndex: -1,
	}

	var err error

	// open tty
	p.tty, err = tty.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open tty: %v", err)
	}

	// listen for SIGINT (Ctrl+C)
	signal.Notify(p.osSignalsChan, syscall.SIGINT)

	// listen for input
	go p.readRunes()

	return &p, nil
}

func (p *Prompt) Close() {
	if p.tty != nil {
		p.tty.Close()
	}
	close(p.osSignalsChan)
	close(p.runeChan)
	close(p.errChan)
}

func (p *Prompt) HandlePrompt(previousInput string) (types.Prompt, string, error) {
	fmt.Print(p.cfg.PromptSymbol + " " + previousInput)

	input, skipExec := p.readInput(previousInput)
	if skipExec {
		return types.Prompt{}, input, nil
	}

	prompt, err := p.parseInput(strings.TrimSpace(input))

	return prompt, "", err
}
