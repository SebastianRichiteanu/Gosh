package builtins

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/SebastianRichiteanu/Gosh/internal/types"
	"github.com/SebastianRichiteanu/Gosh/internal/utils"
)

const (
	BuiltinExit   = "exit"
	BuiltinEcho   = "echo"
	BuiltinPwd    = "pwd"
	BuiltinCd     = "cd"
	BuiltinType   = "type"
	BuiltinClear  = "clear"
	BuiltinSource = "source"
	BuiltinExport = "export"

	ClearControlSeq = "\033[H\033[2J"
)

var builtinCmds types.CommandMap = make(types.CommandMap)
var reloadCfgChannel chan bool = make(chan bool)

// InitBuiltinCmds initializes all built-in commands and stores them in a CommandMap for easy lookup
func InitBuiltinCmds() (types.CommandMap, chan bool) {
	// All fct ret are (string, error) for now at least

	builtinCmds[BuiltinExit] = builtinExit()
	builtinCmds[BuiltinEcho] = builtinEcho()
	builtinCmds[BuiltinPwd] = builtinPwd()
	builtinCmds[BuiltinCd] = builtinCd()
	builtinCmds[BuiltinType] = builtinType()
	builtinCmds[BuiltinClear] = builtinClear()
	builtinCmds[BuiltinSource] = builtinSource()
	builtinCmds[BuiltinExport] = builtinExport()

	return builtinCmds, reloadCfgChannel
}

// builtinExit defines the exit behavior of the shell
// It terminates the program with an optional exit code
func builtinExit() types.Command {
	return func(code ...string) (string, error) {
		if len(code) == 0 {
			utils.ExitShell(0)
		}

		codeAsInt, err := strconv.Atoi(code[0])
		if err != nil {
			return "", err
		}

		utils.ExitShell(codeAsInt)
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
func builtinSource() types.Command {
	return func(filePaths ...string) (string, error) {
		// TODO: double check this function : GOSH_ENABLE_AUTOCOMPLETE=fals
		for _, filePath := range filePaths {
			file, err := os.Open(filePath)
			if err != nil {
				return "", fmt.Errorf("failed to open %s: %w", filePath, err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())

				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}

				// Check if it's a variable assignment (e.g., VAR=value)
				if strings.Contains(line, "=") {
					line = strings.Trim(line, "export")
					line = strings.Trim(line, " ")

					parts := strings.SplitN(line, "=", 2)
					key := strings.TrimSpace(parts[0])
					value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
					os.Setenv(key, value)
					continue
				}

				// Execute command if not
				cmd := exec.Command("sh", "-c", line)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					fmt.Println("Error executing command:", line, err)
				}
			}

			if err := scanner.Err(); err != nil {
				return "", fmt.Errorf("error reading %s: %w", filePath, err)
			}
		}

		reloadCfgChannel <- true

		return "", nil
	}
}

func builtinExport() types.Command {
	return func() {}
	// TODO
}
