package prompt

import (
	"fmt"
	"os"
	"strings"
)

func (p *Prompt) loadAliases() error {
	data, err := os.ReadFile(p.cfg.AliasFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No alias file yet
		}
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		alias := parts[0]
		command := strings.Trim(parts[1], "'\"")

		if p.aliases == nil {
			p.logger.Error("aliases map is nil")
			return nil
		}

		(*p.aliases)[alias] = command
	}

	return nil
}

func (p *Prompt) saveAliases() error {
	f, err := os.OpenFile(p.cfg.AliasFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if p.aliases == nil {
		return nil
	}

	for alias, command := range *p.aliases {
		_, err := f.WriteString(fmt.Sprintf("%s='%s'\n", alias, command))
		if err != nil {
			return err
		}
	}

	return nil
}
