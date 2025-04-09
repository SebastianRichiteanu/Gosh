package types

const (
	Stdout = iota + 1
	Stderr = iota + 1
)

const DefaultStdStream = Stdout
const PathDelimiter = ":"
const PathEnvVar = "PATH"

// Command represents the type for built-in or external commands
type Command interface{}

// CommandMap is a map that stores known commands keyed by their name
type CommandMap map[string]Command

// ParsedPrompt is a structure that holds the tokens from user input, stream redirection details, and other command-related information
type ParsedPrompt struct {
	Tokens       []string
	StdStream    int
	RedirectFile string
	Truncate     bool
}
