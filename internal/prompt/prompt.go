package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/SebastianRichiteanu/Gosh/internal/builtins"
	"github.com/SebastianRichiteanu/Gosh/internal/completer"
	"github.com/SebastianRichiteanu/Gosh/internal/types"
	"github.com/SebastianRichiteanu/Gosh/internal/utils"
)

var (
	errUnterminatedSingleQuote = errors.New("unterminated single quotes")
	errUnterminatedDoubleQuote = errors.New("unterminated double quotes")
)

type Prompt struct {
	knownCmds     *types.CommandMap
	autocompleter *completer.Autocompleter
	history       []string
}

func NewPrompt(knownCmds *types.CommandMap, autocompleter *completer.Autocompleter) *Prompt {
	return &Prompt{
		knownCmds:     knownCmds,
		autocompleter: autocompleter,
		history:       make([]string, 0),
	}
}

// Prompt prints the shell prompt, handles user input, and returns the parsed command and tokens
func (p *Prompt) HandlePrompt(oldInput string) (types.Prompt, string, error) {
	fmt.Fprint(os.Stdout, "$ "+oldInput)

	input, skipExec := p.readInput(oldInput)
	if skipExec {
		return types.Prompt{}, input, nil
	}

	return p.parseInput(strings.TrimSpace(input))
}

// readInput handles reading the user input, processing special characters, and returning the final input string
func (p *Prompt) readInput(oldInput string) (string, bool) {
	input := oldInput
	pressedTab := false

	historyIndex := len(p.history)

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	r := bufio.NewReader(os.Stdin)

loop:
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch c {
		case '\x03': // Ctrl+C
			return builtins.BuiltinExit, false
		case '\x0C': // Ctrl+L
			return builtins.BuiltinClear, false
		case '\r', '\n': // Enter
			fmt.Fprint(os.Stdout, "\r\n")
			if strings.TrimSpace(input) != "" {
				p.history = append(p.history, input)
			}
			break loop
		case '\x7F': // Backspace
			if length := len(input); length > 0 {
				input = input[:length-1]
				fmt.Fprint(os.Stdout, "\b \b")
			}
		case '\t': // Tab
			suffixes, _ := p.autocompleter.Autocomplete(*p.knownCmds, input)
			if len(suffixes) == 0 {
				fmt.Fprintf(os.Stdout, "\a")
				continue
			}

			suffixAppender := " "

			splitInput := input
			if strings.Contains(input, " ") {
				splitInputArr := strings.Split(input, " ")
				if len(splitInputArr) == 0 {
					continue // TODO: not sure?
				}

				// TODO: This entire implementation sucks, need to refactor
				splitInputArr = strings.Split(splitInputArr[len(splitInputArr)-1], "/")

				splitInput = splitInputArr[len(splitInputArr)-1]
				suffixAppender = "/" // TODO: only do this is the file is a dir....
			}

			if len(suffixes) == 1 {
				suffix := suffixes[0]

				input += suffix + suffixAppender
				fmt.Fprint(os.Stdout, suffix+suffixAppender)

				continue
			}

			// 2 or more suffixes
			common := p.autocompleter.FindLongestPrefix(suffixes)
			if common != "" {
				input += common
				fmt.Fprint(os.Stdout, common)
				pressedTab = false
				continue
			}

			if !pressedTab {
				fmt.Fprintf(os.Stdout, "\a")
				pressedTab = true
				continue
			}

			var suffixesWithInput []string
			for _, suffix := range suffixes {
				suffixesWithInput = append(suffixesWithInput, splitInput+suffix)
			}

			fmt.Fprintf(os.Stdout, "\r\n%s\n\r", strings.Join(suffixesWithInput, "  "))
			pressedTab = false

			return input, true // return true so we don't exec
		case '\x1B': // Escape sequences (arrow keys)
			r2, _, _ := r.ReadRune()
			r3, _, _ := r.ReadRune()

			if r2 == '\x5B' { // Arrow key sequence
				if r3 == '\x41' { // Up Arrow
					if historyIndex > 0 {
						historyIndex--
						input = p.history[historyIndex]
						fmt.Print("\r$ " + input) // Overwrite current input
					} else if historyIndex == -1 && len(p.history) > 0 {
						historyIndex = len(p.history) - 1
						input = p.history[historyIndex]
						fmt.Print("\r$ " + input)
					}
				} else if r3 == '\x42' { // Down Arrow
					if historyIndex < len(p.history)-1 {
						historyIndex++
						input = p.history[historyIndex]
						fmt.Print("\r$ " + input)
					} else {
						historyIndex = len(p.history)
						input = ""
						fmt.Print("\r$                                               ") // Clear input line
						// TODO: the above is trivial implementation, need to do prompt handling better, maybe use tty directly
					}
				}
			}

		default:
			input += string(c)
			fmt.Fprint(os.Stdout, string(c))
		}
	}
	return input, false
}

// parseInput parses the user input, breaking it into tokens, handling quotes and escape characters, and detecting redirection
func (p *Prompt) parseInput(input string) (types.Prompt, string, error) {
	parsedPrompt := types.Prompt{
		StdStream: types.DefaultStdStream,
		Truncate:  false,
	}

	var currentToken strings.Builder
	var err error

	inSingleQuote := false
	inDoubleQuote := false
	escaping := false

	for i := 0; i < len(input); i++ {
		char := input[i]

		if escaping {
			currentToken.WriteByte(char)
			escaping = false
			continue
		}

		switch char {
		case '\\':
			if inSingleQuote {
				currentToken.WriteByte(char)
				continue
			}

			if inDoubleQuote {
				if i < len(input)-1 && (input[i+1] == '$' || input[i+1] == '"' || input[i+1] == '\\') {
					escaping = true
					continue
				}
				currentToken.WriteByte(char)
				continue
			}

			escaping = true
		case '\'':
			if inDoubleQuote {
				currentToken.WriteByte(char) // Inside double quotes, treat as a literal
				continue
			}

			inSingleQuote = !inSingleQuote
		case '"':
			if inSingleQuote {
				currentToken.WriteByte(char)
				continue
			}

			inDoubleQuote = !inDoubleQuote
		case ' ':
			if inSingleQuote || inDoubleQuote {
				currentToken.WriteByte(char) // Inside quotes, treat as literal
				continue
			}

			if currentToken.Len() > 0 {
				parsedPrompt.Tokens = append(parsedPrompt.Tokens, currentToken.String()) // Outside quotes, end of token
				currentToken.Reset()
			}
		case '>':
			if inSingleQuote || inDoubleQuote {
				currentToken.WriteByte(char)
				continue
			}

			parsedPrompt.StdStream, err = utils.GetStdStream(input, i-1)
			if err != nil {
				return parsedPrompt, "", err
			}

			if i < len(input)-1 && input[i+1] == '>' {
				parsedPrompt.RedirectFile = strings.Trim(input[i+2:], " \"") // trim spaces and quotes
			} else {
				parsedPrompt.RedirectFile = strings.Trim(input[i+2:], " \"") // trim spaces and quotes
				parsedPrompt.Truncate = true
			}

			return parsedPrompt, "", nil
		default:
			currentToken.WriteByte(char)
		}
	}

	if currentToken.Len() > 0 {
		parsedPrompt.Tokens = append(parsedPrompt.Tokens, currentToken.String())
	}

	if inSingleQuote {
		return parsedPrompt, "", errUnterminatedSingleQuote
	}
	if inDoubleQuote {
		return parsedPrompt, "", errUnterminatedDoubleQuote
	}

	return parsedPrompt, "", nil
}
