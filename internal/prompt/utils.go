package prompt

import (
	"fmt"
	"os"
)

func (p *Prompt) renderPrompt(prompt []rune) {
	fmt.Printf("\r%s %s\033[K", p.cfg.PromptSymbol, string(prompt))
}

func (p *Prompt) bell() {
	fmt.Fprintf(os.Stdout, "\a")
}

func (p *Prompt) moveCursorBack(positions int) {
	if positions <= 0 {
		return
	}

	fmt.Printf("\033[%dD", positions)
}

func (p *Prompt) moveCursorFront(positions int) {
	if positions <= 0 {
		return
	}

	fmt.Printf("\033[%dC", positions)
}
