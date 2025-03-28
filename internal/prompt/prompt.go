package prompt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/SebastianRichiteanu/Gosh/internal/autocompleter"
	"github.com/SebastianRichiteanu/Gosh/internal/logger"
	"github.com/SebastianRichiteanu/Gosh/internal/types"
)

var (
	errUnterminatedSingleQuote = errors.New("unterminated single quotes")
	errUnterminatedDoubleQuote = errors.New("unterminated double quotes")
)

type Prompt struct {
	builtinCmds   *types.CommandMap
	autocompleter *autocompleter.Autocompleter
	logger        *logger.Logger

	history       []string
	historyIndex  int
	promptSymbols string
}

func NewPrompt(builtinCmds *types.CommandMap, autocompleter *autocompleter.Autocompleter, logger *logger.Logger) *Prompt {
	return &Prompt{
		builtinCmds:   builtinCmds,
		autocompleter: autocompleter,
		logger:        logger,
		history:       []string{},
		historyIndex:  -1,
		promptSymbols: "$ ", // TODO: add to config
	}
}

func (p *Prompt) HandlePrompt(previousInput string) (types.Prompt, string, error) {
	fmt.Print(p.promptSymbols + previousInput)

	input, skipExec := p.readInput(previousInput)
	if skipExec {
		return types.Prompt{}, input, nil
	}

	return p.parseInput(strings.TrimSpace(input))
}
