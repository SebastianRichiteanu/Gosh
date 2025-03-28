package prompt

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/SebastianRichiteanu/Gosh/internal/builtins"
	"github.com/SebastianRichiteanu/Gosh/internal/completer"
	"github.com/SebastianRichiteanu/Gosh/internal/types"
	"github.com/SebastianRichiteanu/Gosh/internal/utils"
	"github.com/mattn/go-tty"
)

var (
	errUnterminatedSingleQuote = errors.New("unterminated single quotes")
	errUnterminatedDoubleQuote = errors.New("unterminated double quotes")
)

type Prompt struct {
	knownCmds     *types.CommandMap
	autocompleter *completer.Autocompleter
	history       []string
	historyIndex  int
}

func NewPrompt(knownCmds *types.CommandMap, autocompleter *completer.Autocompleter) *Prompt {
	return &Prompt{
		knownCmds:     knownCmds,
		autocompleter: autocompleter,
		history:       []string{},
		historyIndex:  -1,
	}
}

func (p *Prompt) HandlePrompt(previousInput string) (types.Prompt, string, error) {
	fmt.Print("$ " + previousInput)
	input, skipExec := p.readInput(previousInput)
	if skipExec {
		return types.Prompt{}, input, nil
	}
	return p.parseInput(strings.TrimSpace(input))
}

func (p *Prompt) readInput(previousInput string) (string, bool) {
	tty, err := tty.Open()
	if err != nil {
		panic(err)
	}
	defer tty.Close()

	input := []rune(previousInput)
	cursor := len(previousInput)
	pressedTab := false

	for {
		char, err := tty.ReadRune()
		if err != nil {
			continue
		}

		switch char {
		case 12: // Ctrl+L
			return builtins.BuiltinClear, false
		case 13: // Enter
			fmt.Println()
			if len(input) > 0 {
				p.history = append(p.history, string(input))
			}
			p.historyIndex = len(p.history)
			return string(input), false
		case 127: // Backspace
			if cursor > 0 {
				input = append(input[:cursor-1], input[cursor:]...)
				cursor--
				fmt.Printf("\r$ %s \033[K", string(input))  // Clear line after cursor
				fmt.Printf("\033[%dD", len(input)-cursor+1) // Move cursor back
			}
		case 9: // Tab (Autocomplete)
			// TODO: move the below and maybe only handle runes?

			inputAsStr := string(input)

			suffixes, _ := p.autocompleter.Autocomplete(*p.knownCmds, inputAsStr)
			if len(suffixes) == 0 {
				fmt.Fprintf(os.Stdout, "\a")
				continue
			}

			suffixAppender := " "

			splitInput := inputAsStr
			if strings.Contains(inputAsStr, " ") {
				splitInputArr := strings.Split(inputAsStr, " ")
				if len(splitInputArr) == 0 {
					continue // TODO: not sure?
				}

				splitInputArr = strings.Split(splitInputArr[len(splitInputArr)-1], "/")

				splitInput = splitInputArr[len(splitInputArr)-1]
				suffixAppender = "/" // TODO: only do this is the file is a dir....
			}

			if len(suffixes) == 1 {
				suffix := suffixes[0]

				input = append(input, []rune(suffix)...)
				input = append(input, []rune(suffixAppender)...)

				fmt.Fprint(os.Stdout, suffix+suffixAppender)

				continue
			}

			// 2 or more suffixes
			common := p.autocompleter.FindLongestPrefix(suffixes)
			if common != "" {
				input = append(input, []rune(common)...)
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

			return string(input), true // return true so we don't exec
		case 27: // Escape sequences (Arrow keys)
			if r2, _ := tty.ReadRune(); r2 == 91 {
				switch r3, _ := tty.ReadRune(); r3 {
				case 65: // Up Arrow
					if p.historyIndex > 0 {
						p.historyIndex--
						input = []rune(p.history[p.historyIndex])
						cursor = len(input)
						fmt.Printf("\r$ %s\033[K", string(input))
					}
				case 66: // Down Arrow
					if p.historyIndex < len(p.history)-1 {
						p.historyIndex++
						input = []rune(p.history[p.historyIndex])
						cursor = len(input)
						fmt.Printf("\r$ %s\033[K", string(input))
					} else {
						p.historyIndex = len(p.history)
						input = []rune{}
						fmt.Print("\r$ \033[K")
					}
				case 67: // Right Arrow
					if cursor < len(input) {
						cursor++
						fmt.Print("\033[C")
					}
				case 68: // Left Arrow
					if cursor > 0 {
						cursor--
						fmt.Print("\033[D")
					}
				}
			}
		default:
			input = append(input[:cursor], append([]rune{char}, input[cursor:]...)...)
			cursor++
			fmt.Printf("\r$ %s ", string(input))
			fmt.Printf("\033[%dD", len(input)-cursor+1) // Move cursor back to correct position
		}
	}
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
