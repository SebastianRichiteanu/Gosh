package executor

import (
	"fmt"
	"reflect"

	"github.com/SebastianRichiteanu/Gosh/internal/types"
)

// execBuiltin executes a built-in command and handles its output (stdout, stderr)
func (e *Executor) execBuiltin(knownCmd types.Command, prompt types.ParsedPrompt) {
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

// runBuiltin runs a built-in function dynamically using reflection, passing the arguments from the prompt
func (e *Executor) runBuiltin(cmd types.Command, prompt types.ParsedPrompt) []reflect.Value {
	args := prompt.Tokens[1:]

	fctValue := reflect.ValueOf(cmd)
	fctType := fctValue.Type()

	numIn := fctType.NumIn()
	isVariadic := fctType.IsVariadic()

	if (!isVariadic && numIn != len(args)) || (isVariadic && len(args) < numIn-1) {
		expected := fmt.Sprintf("%d", numIn)
		if isVariadic {
			expected = fmt.Sprintf("at least %d", numIn-1)
		}
		errMsg := fmt.Sprintf("error: wrong number of arguments (expected %s, got %d)", expected, len(args))
		return []reflect.Value{
			reflect.ValueOf(""),
			reflect.ValueOf(errMsg),
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
