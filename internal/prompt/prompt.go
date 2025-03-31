package prompt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/SebastianRichiteanu/Gosh/internal/autocompleter"
	"github.com/SebastianRichiteanu/Gosh/internal/config"
	"github.com/SebastianRichiteanu/Gosh/internal/logger"
	"github.com/SebastianRichiteanu/Gosh/internal/types"
)

var (
	errUnterminatedSingleQuote = errors.New("unterminated single quotes")
	errUnterminatedDoubleQuote = errors.New("unterminated double quotes")
)

type Prompt struct {
	cfg           config.Config
	builtinCmds   *types.CommandMap
	autocompleter *autocompleter.Autocompleter
	logger        *logger.Logger

	history      []string
	historyIndex int
}

func NewPrompt(builtinCmds *types.CommandMap, autocompleter *autocompleter.Autocompleter, cfg *config.Config, logger *logger.Logger) *Prompt {
	return &Prompt{
		cfg:           *cfg,
		builtinCmds:   builtinCmds,
		autocompleter: autocompleter,
		logger:        logger,
		history:       []string{},
		historyIndex:  -1,
	}
}

func (p *Prompt) HandlePrompt(previousInput string) (types.Prompt, string, error) {
	p.logger.Info(fmt.Sprintf("-----%#v -----", p.cfg)) // TODO: this is not getting updated.....

	fmt.Print(p.cfg.PromptSymbol + " " + previousInput)

	input, skipExec := p.readInput(previousInput)
	if skipExec {
		return types.Prompt{}, input, nil
	}

	return p.parseInput(strings.TrimSpace(input))
}
