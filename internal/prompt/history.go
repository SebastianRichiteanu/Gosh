package prompt

import (
	"os"
	"strings"
)

func (p *Prompt) loadHistory() error {
	data, err := os.ReadFile(p.cfg.HistoryFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // no file yet
		}
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			p.history = append(p.history, line)
		}
	}

	p.historyIndex = len(p.history)
	return nil
}

func (p *Prompt) appendToHistoryFile(cmd string) error {
	f, err := os.OpenFile(p.cfg.HistoryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(cmd + "\n")
	return err
}

func (p *Prompt) rewriteHistoryFile(history []string) error {
	f, err := os.OpenFile(p.cfg.HistoryFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range history {
		if _, err := f.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}
