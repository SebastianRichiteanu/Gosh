package builtins

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SebastianRichiteanu/Gosh/internal/types"
	"github.com/SebastianRichiteanu/Gosh/internal/utils"
)

const (
	BuiltinExit    = "exit"
	BuiltinEcho    = "echo"
	BuiltinPwd     = "pwd"
	BuiltinCd      = "cd"
	BuiltinType    = "type"
	BuiltinClear   = "clear"
	BuiltinSource  = "source"
	BuiltinExport  = "export"
	BuiltinHistory = "history"

	ClearControlSeq = "\033[H\033[2J"
)

var builtinCmds types.CommandMap = make(types.CommandMap)

// InitBuiltinCmds initializes all built-in commands and stores them in a CommandMap for easy lookup
func InitBuiltinCmds(exitChannel chan int, reloadCfgChannel chan bool, historyFile *string) types.CommandMap {
	// All fct ret are (string, error) for now at least

	builtinCmds[BuiltinExit] = builtinExit(exitChannel)
	builtinCmds[BuiltinEcho] = builtinEcho()
	builtinCmds[BuiltinPwd] = builtinPwd()
	builtinCmds[BuiltinCd] = builtinCd()
	builtinCmds[BuiltinType] = builtinType()
	builtinCmds[BuiltinClear] = builtinClear()
	builtinCmds[BuiltinSource] = builtinSource(reloadCfgChannel)
	builtinCmds[BuiltinExport] = builtinExport(reloadCfgChannel)
	builtinCmds[BuiltinHistory] = builtinHistory(historyFile)

	return builtinCmds
}

// builtinExit defines the exit behavior of the shell
// It terminates the program with an optional exit code
func builtinExit(exitChannel chan int) types.Command {
	return func(code ...string) (string, error) {
		if len(code) == 0 {
			exitChannel <- 0
			return "", nil
		}

		codeAsInt, err := strconv.Atoi(code[0])
		if err != nil {
			return "", err
		}

		exitChannel <- codeAsInt
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
		if _, isKnownCmd := builtinCmds[cmd]; isKnownCmd || cmd == BuiltinType {
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

// builtinSource defines the source behavior of the shell
// It will source a file
func builtinSource(reloadCfgChannel chan bool) types.Command {
	return func(filePaths ...string) (string, error) {
		for _, filePath := range filePaths {
			if err := utils.SourceFile(filePath); err != nil {
				return "", fmt.Errorf("could not source file at path %s: %v", filePath, err)
			}
		}

		reloadCfgChannel <- true

		time.Sleep(time.Millisecond) // The update is a bit slow so wait a milisec before returning

		return "", nil
	}
}

func builtinExport(reloadCfgChannel chan bool) types.Command {
	return func(line string) (string, error) {
		utils.HandleExportLine(line)
		reloadCfgChannel <- true

		time.Sleep(time.Millisecond) // The update is a bit slow so wait a milisec before returning

		return "", nil
	}
}

func builtinHistory(historyFile *string) types.Command {
	return func(args ...string) (string, error) {
		if historyFile == nil {
			return "", fmt.Errorf("history file not set")
		}

		data, err := os.ReadFile(*historyFile)
		if err != nil {
			if os.IsNotExist(err) {
				return "", nil // no file yet
			}
			return "", err
		}

		var sb strings.Builder

		lines := strings.Split(string(data), "\n")
		for idx, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				sb.WriteString(fmt.Sprintf("%d  %s\n", idx+1, line))
			}
		}

		return sb.String(), nil
	}
}
