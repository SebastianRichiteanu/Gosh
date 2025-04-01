package prompt

import (
	"fmt"
	"os"
	"strings"

	"github.com/SebastianRichiteanu/Gosh/internal/builtins"
	"github.com/SebastianRichiteanu/Gosh/internal/utils"
	"github.com/mattn/go-tty"
)

func (p *Prompt) renderPrompt(prompt []rune) {
	fmt.Printf("\r%s %s\033[K", p.cfg.PromptSymbol, string(prompt))
}

func (p *Prompt) bell() {
	fmt.Fprintf(os.Stdout, "\a")
}

func (p *Prompt) moveCursorBack(positions int) {
	fmt.Printf("\033[%dD", positions)
}

func (p *Prompt) moveCursorFront(positions int) {
	fmt.Printf("\033[%dC", positions)
}

func (p *Prompt) readInput(previousInput string) (string, bool) {
	tty, err := tty.Open()
	if err != nil {
		panic(err)
	}
	defer tty.Close()

	input := []rune(previousInput)
	inputBkp := []rune{}

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
				p.renderPrompt(append(input, ' ')) // idk why tbh, but it works
				p.moveCursorBack(len(input) - cursor + 1)
			}
		case 9: // Tab (Autocomplete)
			// TODO: move the below and maybe only handle runes?
			inputAsStr := string(input)

			currentPrompt, err := p.parseInput(inputAsStr)
			if err != nil {
				p.logger.Error(fmt.Sprintf("failed to parse input: %v", err), "input", inputAsStr)
				p.bell()
				continue
			}

			tokenIndex := p.findTokenIndexAtPosition(currentPrompt.Tokens, cursor)
			tokenToAutocomplete := currentPrompt.Tokens[tokenIndex]
			suffixes := p.autocompleter.Autocomplete(*p.builtinCmds, tokenToAutocomplete)
			if len(suffixes) == 0 {
				p.bell()
				continue
			}

			if len(suffixes) == 1 {
				// add the suffix in the prompt
				// TODO: this will only work if we are on the last char of the token, right?
				// Have to do the below with rerender prompt as well

				suffix := suffixes[0]
				if suffix[len(suffix)-1] != '/' {
					suffix += " "
				}

				//input = append(append(input[:cursor], []rune(suffix)...), input[cursor:]...)
				input = append(input[:cursor], append([]rune(suffix), input[cursor:]...)...)

				cursor += len(suffix)
				p.renderPrompt(input)

				continue
			}

			// 2 or more suffixes
			common := utils.FindLongestPrefix(suffixes)
			if common != "" {
				input = append(input, []rune(common)...)
				cursor += len(common)
				pressedTab = false

				p.renderPrompt(input)

				continue
			}

			if !pressedTab {
				p.bell()
				pressedTab = true
				continue
			}

			pathToken := strings.Split(tokenToAutocomplete, "/")
			tokenToAutocomplete = pathToken[len(pathToken)-1]

			var suffixesWithInput []string
			for _, suffix := range suffixes {
				suffixesWithInput = append(suffixesWithInput, tokenToAutocomplete+suffix)
			}

			fmt.Fprintf(os.Stdout, "\r\n%s\n\r", strings.Join(suffixesWithInput, "  "))
			pressedTab = false

			return string(input), true // return true so we don't exec
		case 27: // Escape sequences (Arrow keys)
			if r2, _ := tty.ReadRune(); r2 == 91 {
				switch r3, _ := tty.ReadRune(); r3 {
				case 65: // Up Arrow
					if p.historyIndex > 0 {
						if p.historyIndex == len(p.history) {
							inputBkp = input
						}

						p.historyIndex--
						input = []rune(p.history[p.historyIndex])
						cursor = len(input)
						p.renderPrompt(input)
					}
				case 66: // Down Arrow
					if p.historyIndex < len(p.history)-1 { // TODO: I think the historyIndex is broken
						p.historyIndex++
						input = []rune(p.history[p.historyIndex])
						cursor = len(input)
						p.renderPrompt(input)
					} else if p.historyIndex == len(p.history)-1 {
						p.historyIndex = len(p.history)
						input = inputBkp
						cursor = len(input)
						p.renderPrompt(input)
						inputBkp = []rune{}
					} else {
						p.bell()
						continue
					}
				case 67: // Right Arrow
					if cursor < len(input) {
						cursor++
						p.moveCursorFront(1)
					}
				case 68: // Left Arrow
					if cursor > 0 {
						cursor--
						p.moveCursorBack(1)
					}
				}
			}
		default:
			input = append(input[:cursor], append([]rune{char}, input[cursor:]...)...)
			cursor++
			fmt.Printf("\r%s %s ", p.cfg.PromptSymbol, string(input))
			fmt.Printf("\033[%dD", len(input)-cursor+1) // Move cursor back to correct position
		}
	}
}
