package executor

import (
	"github.com/SebastianRichiteanu/Gosh/internal/types"
)

type Executor struct {
	builtinCmds *types.CommandMap
}

func NewExecutor(builtinCmds *types.CommandMap) *Executor {
	return &Executor{
		builtinCmds: builtinCmds,
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
