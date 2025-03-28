package executor

import (
	"github.com/SebastianRichiteanu/Gosh/internal/logger"
	"github.com/SebastianRichiteanu/Gosh/internal/types"
)

type Executor struct {
	builtinCmds *types.CommandMap
	logger      *logger.Logger
}

func NewExecutor(builtinCmds *types.CommandMap, logger *logger.Logger) *Executor {
	return &Executor{
		builtinCmds: builtinCmds,
		logger:      logger,
	}
}

// Execute executes the given command based on the parsed prompt
func (e *Executor) Execute(prompt types.Prompt) {
	if len(prompt.Tokens) == 0 {
		return
	}

	if e.builtinCmds != nil {
		knownCmd, isKnownCmd := (*e.builtinCmds)[prompt.Tokens[0]]
		if isKnownCmd {
			e.execBuiltin(knownCmd, prompt)
			return
		}
	}

	e.execBinary(prompt)
}
