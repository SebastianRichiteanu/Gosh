package prompt

import (
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
