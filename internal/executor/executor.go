package executor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"

	"github.com/SebastianRichiteanu/Gosh/internal/types"
	"github.com/SebastianRichiteanu/Gosh/internal/utils"
)

type Executor struct {
	knownCmds *types.CommandMap
}

func NewExecutor(knownCmds *types.CommandMap) *Executor {
	return &Executor{
		knownCmds: knownCmds,
	}
}

// Execute executes the given command based on the parsed prompt
func (e *Executor) Execute(prompt types.Prompt) {
	if len(prompt.Tokens) == 0 {
		return
	}

	if e.knownCmds != nil {
		knownCmd, isKnownCmd := (*e.knownCmds)[prompt.Tokens[0]]
		if isKnownCmd {
			e.execBuiltin(knownCmd, prompt)
			return
		}
	}

	e.execBinary(prompt)
}

// execBuiltin executes a built-in command and handles its output (stdout, stderr)
func (e *Executor) execBuiltin(knownCmd types.Command, prompt types.Prompt) {
	output := e.runBuiltin(knownCmd, prompt)

	if len(output) != 2 {
		panic(fmt.Errorf("command did not return 2 out streams"))
	}

	stdout := output[0]
	stderr := output[1]

	if prompt.RedirectFile == "" {
		e.handleDirectOutput(stdout, stderr)
		return
	}

	e.handleFileOutput(prompt, stdout, stderr)
}

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
		file.WriteString(stdout.String())
		if !stderr.IsNil() {
			fmt.Println(stderr)
		}
	case types.Stderr:
		fmt.Print(stdout)
		if !stderr.IsNil() {
			file.WriteString(stderr.String())
		}
	default:
		panic(errors.New("wtf?"))
	}
}

// runBuiltin runs a built-in function dynamically using reflection, passing the arguments from the prompt
func (e *Executor) runBuiltin(cmd types.Command, prompt types.Prompt) []reflect.Value {
	args := prompt.Tokens[1:]

	fctValue := reflect.ValueOf(cmd)
	fctType := fctValue.Type()

	numIn := fctType.NumIn()
	isVariadic := fctType.IsVariadic()

	if (!isVariadic && numIn != len(args)) || (isVariadic && len(args) < numIn-1) {
		// Return 2 streams, first for stdout and then for stderr
		return []reflect.Value{
			reflect.ValueOf(""),
			reflect.ValueOf(fmt.Sprintf("%s is not a callable function", cmd)),
		}
	}

	fctArgs := make([]reflect.Value, 0, numIn)

	for i := 0; i < numIn; i++ {
		if isVariadic && i == numIn-1 {
			variadicSlice := reflect.MakeSlice(fctType.In(i), len(args[i:]), len(args[i:]))
			for j, arg := range args[i:] {
				variadicSlice.Index(j).Set(reflect.ValueOf(arg))
			}
			fctArgs = append(fctArgs, variadicSlice)
			continue
		}

		fctArgs = append(fctArgs, reflect.ValueOf(args[i]))
	}

	if isVariadic {
		return fctValue.CallSlice(fctArgs)
	}

	return fctValue.Call(fctArgs)
}

// execBinary executes an external command by searching for the binary in the system's PATH
// and invoking it with the provided arguments
func (e *Executor) execBinary(prompt types.Prompt) {
	binary := prompt.Tokens[0]
	args := prompt.Tokens[1:]

	fullPath := utils.FindPath(binary)
	if fullPath == "" {
		fmt.Printf("%s: not found\n", binary)
		return
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
			panic(errors.New("how?"))
		}
	}

	cmd.Run()
}
