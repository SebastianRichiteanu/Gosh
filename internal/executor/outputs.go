package executor

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/SebastianRichiteanu/Gosh/internal/types"
	"github.com/SebastianRichiteanu/Gosh/internal/utils"
)

// handleDirectOutput prints the command's stdout or stderr directly to the terminal
func (e *Executor) handleDirectOutput(stdout, stderr reflect.Value) {
	if stderr.IsNil() {
		fmt.Printf("%s", stdout.String())
		return
	}
	val := stderr.Interface()
	fmt.Printf("%s\n", val.(error))
}

// handleFileOutput handles redirecting command output to a file based on the prompt's redirection
func (e *Executor) handleFileOutput(prompt types.Prompt, stdout, stderr reflect.Value) {
	file, err := utils.OpenFileForStdout(prompt.RedirectFile, prompt.Truncate)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer file.Close()

	switch prompt.StdStream {
	case types.Stdout:
		if _, err := file.WriteString(stdout.String()); err != nil {
			e.logger.Error(fmt.Sprintf("failed to write string for stdout: %v", err))
		}
		if !stderr.IsNil() {
			fmt.Println(stderr)
		}
	case types.Stderr:
		fmt.Print(stdout)
		if !stderr.IsNil() {
			if _, err := file.WriteString(stderr.String()); err != nil {
				e.logger.Error(fmt.Sprintf("failed to write string for stderr: %v", err))
			}
		}
	default:
		panic(errors.New("wtf?"))
	}
}
