package prompt

import (
	"fmt"
	"os"
	"strings"

	"github.com/SebastianRichiteanu/Gosh/internal/builtins"
	"github.com/SebastianRichiteanu/Gosh/internal/utils"
)

func (p *Prompt) readRunes() {
	for {
		r, err := p.tty.ReadRune()
		if err != nil {
			p.errChan <- err
			return
		}

		// If it's an escape character, read more to handle arrow keys
		if r == runeEscape {
			r2, err := p.tty.ReadRune()
			if err != nil {
				p.errChan <- err
				return
			}
			if r2 != runeBracket { // Not a control sequence
				p.runeChan <- r
				p.runeChan <- r2
				continue
			}

			r3, err := p.tty.ReadRune()
			if err != nil {
				p.errChan <- err
				return
			}

			switch r3 {
			case runeArrowUp:
				p.runeChan <- myRuneArrowUp
			case runeArrowDown:
				p.runeChan <- myRuneArrowDown
			case runeArrowRight:
				p.runeChan <- myRuneArrowRight
			case runeArrowLeft:
				p.runeChan <- myRuneArrowLeft
			default:
				// Unknown sequence, send all
				p.runeChan <- r
				p.runeChan <- r2
				p.runeChan <- r3
			}
			continue
		}

		p.runeChan <- r
	}
}

func (p *Prompt) readInput(previousInput string) (string, bool) {
	input := []rune(previousInput)
	inputBkp := []rune{}

	cursor := len(previousInput)
	pressedTab := false

	for {
		select {
		case <-p.osSignalsChan:
			fmt.Println("^C")
			return "", false
		case err := <-p.errChan:
			p.logger.Error(fmt.Sprintf("read error for rune: %v", err))
			continue
		case char := <-p.runeChan:
			// TODO: move each case to a function

			switch char {
			case runeCtrlL:
				return builtins.BuiltinClear, false
			case runeEnter:
				fmt.Println()
				if len(input) > 0 {
					p.history = append(p.history, string(input))
				}
				p.historyIndex = len(p.history)
				return string(input), false
			case runeBackspace:
				if cursor > 0 {
					input = append(input[:cursor-1], input[cursor:]...)
					cursor--
					p.renderPrompt(append(input, ' ')) // idk why tbh, but it works
					p.moveCursorBack(len(input) - cursor + 1)
				}
			case runeTab:
				// TODO: move the below and maybe only handle runes?
				inputAsStr := string(input)

				currentPrompt, err := p.parseInput(inputAsStr)
				if err != nil {
					p.logger.Error(fmt.Sprintf("failed to parse input: %v", err), "input", inputAsStr)
					p.bell()
					continue
				}

				tokenIndex := p.findTokenIndexAtPosition(currentPrompt.Tokens, cursor-1)
				tokenToAutocomplete := currentPrompt.Tokens[tokenIndex]
				suffixes := p.autocompleter.Autocomplete(*p.builtinCmds, tokenToAutocomplete)
				if len(suffixes) == 0 {
					p.bell()
					continue
				}

				if len(suffixes) == 1 {
					// add the suffix in the prompt

					suffix := suffixes[0]
					if len(suffix) > 0 && suffix[len(suffix)-1] != '/' {
						suffix += " "
					}

					// check if we are inside token

					cursorAfterToken := cursor
					for cursorAfterToken < len(input) {
						currentCursorChar := input[cursorAfterToken]
						if currentCursorChar == ' ' {
							break
						}
						cursorAfterToken++
					}

					input = append(input[:cursorAfterToken], append([]rune(suffix), input[cursorAfterToken:]...)...)

					difForTokenEnd := cursorAfterToken - cursor

					cursor += len(suffix) + difForTokenEnd
					p.renderPrompt(input)

					p.moveCursorBack(len(input) - cursor)

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

				if tokenToAutocomplete[len(tokenToAutocomplete)-1] == '/' {
					tokenToAutocomplete = ""
				} else {
					pathToken := strings.Split(tokenToAutocomplete, "/")
					tokenToAutocomplete = pathToken[len(pathToken)-1]
				}

				var suffixesWithInput []string
				for _, suffix := range suffixes {
					suffixesWithInput = append(suffixesWithInput, tokenToAutocomplete+suffix)
				}

				fmt.Fprintf(os.Stdout, "\r\n%s\n\r", strings.Join(suffixesWithInput, "  "))
				pressedTab = false

				return string(input), true // return true so we don't exec
			case myRuneArrowUp:
				if p.historyIndex > 0 {
					if p.historyIndex == len(p.history) {
						inputBkp = input
					}

					p.historyIndex--
					input = []rune(p.history[p.historyIndex])
					cursor = len(input)
					p.renderPrompt(input)
				}
			case myRuneArrowDown:
				if p.historyIndex < len(p.history)-1 {
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
			case myRuneArrowRight:
				if cursor < len(input) {
					cursor++
					p.moveCursorFront(1)
				}
			case myRuneArrowLeft:
				if cursor > 0 {
					cursor--
					p.moveCursorBack(1)
				}
			default:
				input = append(input[:cursor], append([]rune{char}, input[cursor:]...)...)
				cursor++
				fmt.Printf("\r%s %s ", p.cfg.PromptSymbol, string(input))
				fmt.Printf("\033[%dD", len(input)-cursor+1) // Move cursor back to correct position
			}
		}
	}
}
