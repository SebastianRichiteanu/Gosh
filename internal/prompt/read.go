package prompt

import (
	"fmt"
	"os"
	"strings"

	"github.com/SebastianRichiteanu/Gosh/internal/builtins"
	"github.com/SebastianRichiteanu/Gosh/internal/utils"
	"github.com/mattn/go-tty"
)

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
				fmt.Printf("\r$ %s \033[K", string(input))  // Clear line after cursor
				fmt.Printf("\033[%dD", len(input)-cursor+1) // Move cursor back
			}
		case 9: // Tab (Autocomplete)
			// TODO: move the below and maybe only handle runes?

			inputAsStr := string(input)

			suffixes, _ := p.autocompleter.Autocomplete(*p.builtinCmds, inputAsStr)
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
			common := utils.FindLongestPrefix(suffixes)
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
						if p.historyIndex == len(p.history) {
							inputBkp = input
						}

						p.historyIndex--
						input = []rune(p.history[p.historyIndex])
						cursor = len(input)
						fmt.Printf("\r$ %s\033[K", string(input))
					}
				case 66: // Down Arrow
					if p.historyIndex < len(p.history)-1 { // TODO: I think the historyIndex is broken
						p.historyIndex++
						input = []rune(p.history[p.historyIndex])
						cursor = len(input)
						fmt.Printf("\r$ %s\033[K", string(input))
					} else if p.historyIndex == len(p.history)-1 {
						p.historyIndex = len(p.history)
						input = inputBkp
						cursor = len(input)
						fmt.Print("\r$ \033[K", string(input))
						inputBkp = []rune{}
					} else {
						fmt.Fprintf(os.Stdout, "\a")
						continue
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
