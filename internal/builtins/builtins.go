package builtins

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/SebastianRichiteanu/Gosh/internal/types"
	"github.com/SebastianRichiteanu/Gosh/internal/utils"
)

const (
	BuiltinExit  = "exit"
	BuiltinEcho  = "echo"
	BuiltinPwd   = "pwd"
	BuiltinCd    = "cd"
	BuiltinType  = "type"
	BuiltinClear = "clear"

	ClearControlSeq = "\033[H\033[2J"
)

var knownCmds types.CommandMap = make(types.CommandMap)

// InitBuiltins initializes all built-in commands and stores them in a CommandMap for easy lookup
func InitBuiltins() types.CommandMap {
	// All fct ret are (string, error) for now at least

	knownCmds[BuiltinExit] = builtinExit()
	knownCmds[BuiltinEcho] = builtinEcho()
	knownCmds[BuiltinPwd] = builtinPwd()
	knownCmds[BuiltinCd] = builtinCd()
	knownCmds[BuiltinType] = builtinType()
	knownCmds[BuiltinClear] = builtinClear()

	return knownCmds
}

// builtinExit defines the exit behavior of the shell
// It terminates the program with an optional exit code
func builtinExit() types.Command {
	return func(code ...string) (string, error) {
		if len(code) == 0 {
			os.Exit(0)
		}

		codeAsInt, err := strconv.Atoi(code[0])
		if err != nil {
			return "", err
		}

		os.Exit(codeAsInt)
		return "", nil
	}
}

// builtinEcho defines the echo behavior of the shell
// It prints the provided arguments to stdout
func builtinEcho() types.Command {
	return func(args ...string) (string, error) {
		return strings.Join(args, " ") + "\n", nil
	}
}

// builtinPwd defines the pwd behavior of the shell
// It prints the current working directory
func builtinPwd() types.Command {
	return func() (string, error) {
		currentDir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return currentDir + "\n", nil
	}
}

// builtinCd defines the cd behavior of the shell
// It changes the current working directory to the given path
func builtinCd() types.Command {
	return func(dir string) (string, error) {
		expandedPath, err := utils.ExpandHomePath(dir)
		if err != nil {
			return "", err
		}

		if err := os.Chdir(expandedPath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				dir = strings.TrimPrefix(dir, "/")
				return "", fmt.Errorf("%s: /%s: No such file or directory", BuiltinCd, dir)
			}
			return "", err
		}

		return "", nil
	}
}

// builtinType defines the type behavior of the shell
// It prints the type of a given command (either a built-in or external command).
func builtinType() types.Command {
	return func(cmd string) (string, error) {
		if _, isKnownCmd := knownCmds[cmd]; isKnownCmd || cmd == BuiltinType {
			return "", fmt.Errorf("%s is a shell builtin", cmd)
		}

		fullPath := utils.FindPath(cmd)
		if fullPath == "" {
			return "", fmt.Errorf("%s: not found", cmd)
		}

		return fmt.Sprintf("%s is %s\n", cmd, fullPath), nil
	}
}

// builtinClear defines the clear behavior of the shell
// It clears the terminal screen.
func builtinClear() types.Command {
	return func() (string, error) {
		return ClearControlSeq, nil
	}
}
