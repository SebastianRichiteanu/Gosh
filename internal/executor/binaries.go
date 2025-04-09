package executor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/SebastianRichiteanu/Gosh/internal/types"
	"github.com/SebastianRichiteanu/Gosh/internal/utils"
)

// execBinary executes an external command by searching for the binary in the system's PATH
// and invoking it with the provided arguments
func (e *Executor) execBinary(prompt types.ParsedPrompt) {
	binary := prompt.Tokens[0]
	args := prompt.Tokens[1:]

	fullPath := utils.FindPath(binary)
	if fullPath == "" {
		fmt.Printf("%s: not found\n", binary)
		return
	}

	for idx, arg := range args {
		if arg[0] == '~' {
			fullArg, err := utils.ExpandHomePath(arg)
			if err != nil {
				continue
			}

			args[idx] = fullArg
		}
	}

	cmd := exec.Command(binary, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if prompt.RedirectFile != "" {
		outFile, err := utils.OpenFileForStdout(prompt.RedirectFile, prompt.Truncate)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		switch prompt.StdStream {
		case types.Stdout:
			cmd.Stdout = outFile
		case types.Stderr:
			cmd.Stderr = outFile
		default:
			panic(errors.New("unable to handle unknown std stream"))
		}
	}

	cmd.Run()
}
