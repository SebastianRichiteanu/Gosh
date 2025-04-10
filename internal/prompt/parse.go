package prompt

import (
	"os"
	"strings"

	"github.com/SebastianRichiteanu/Gosh/internal/types"
	"github.com/SebastianRichiteanu/Gosh/internal/utils"
)

func (p *Prompt) findTokenIndexAtPosition(tokens []string, position int) int {
	index := 0
	for idx, token := range tokens {
		index += len(token)

		if index >= position {
			return idx
		}
	}

	return len(tokens) - 1
}

// parseInput parses the user input, breaking it into tokens, handling quotes and escape characters, and detecting redirection
func (p *Prompt) parseInput(input string) (types.ParsedPrompt, error) {
	parsedPrompt := types.ParsedPrompt{
		StdStream: types.DefaultStdStream,
		Truncate:  false,
	}

	var currentToken strings.Builder
	var err error

	tokens := strings.Fields(input)
	if len(tokens) > 0 && p.aliases != nil {
		if aliasCommand, exists := (*p.aliases)[tokens[0]]; exists {
			expandedInput := aliasCommand + " " + strings.Join(tokens[1:], " ")
			return p.parseInput(expandedInput)
		}
	}

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
				return parsedPrompt, err
			}

			if i < len(input)-1 && input[i+1] == '>' {
				parsedPrompt.RedirectFile = strings.Trim(input[i+2:], " \"") // trim spaces and quotes
			} else {
				parsedPrompt.RedirectFile = strings.Trim(input[i+2:], " \"") // trim spaces and quotes
				parsedPrompt.Truncate = true
			}

			return parsedPrompt, nil
		case '$':
			if inSingleQuote || inDoubleQuote {
				currentToken.WriteByte(char)
				continue
			}

			varName := ""
			if i+1 < len(input) && input[i+1] == '{' {
				// ${VAR_NAME}
				end := strings.IndexByte(input[i+2:], '}')
				if end != -1 {
					varName = input[i+2 : i+2+end]
					i += end + 2
				} else {
					currentToken.WriteByte(char) // End not found, treat as literal
					continue
				}
			} else {
				// $VAR_NAME
				j := i + 1
				for j < len(input) && isAlphaNumeric(input[j]) {
					j++
				}
				varName = input[i+1 : j]
				i = j - 1
			}

			if val := os.Getenv(varName); val != "" {
				currentToken.WriteString(val)
			}

		default:
			currentToken.WriteByte(char)
		}
	}

	if currentToken.Len() > 0 {
		parsedPrompt.Tokens = append(parsedPrompt.Tokens, currentToken.String())
	}

	if inSingleQuote {
		return parsedPrompt, errUnterminatedSingleQuote
	}
	if inDoubleQuote {
		return parsedPrompt, errUnterminatedDoubleQuote
	}

	return parsedPrompt, nil
}

func isAlphaNumeric(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9') ||
		b == '_'
}
