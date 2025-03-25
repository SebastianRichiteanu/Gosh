package utils

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/SebastianRichiteanu/Gosh/internal/types"
)

// FindPath searches the system's PATH for a given command and returns its full file path if found
func FindPath(cmd string) string {
	paths := strings.Split(os.Getenv(types.PathEnvVar), types.PathDelimiter)
	for _, path := range paths {
		fp := filepath.Join(path, cmd)
		if _, err := os.Stat(fp); err == nil {
			return fp
		}
	}
	return ""
}

// ExpandHomePath expands `~` to the user's home directory in the provided path
func ExpandHomePath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, path[1:]), nil
	}
	return path, nil
}

// OpenFileForStdout opens a file for output redirection, ensuring that the parent directories exist, and returns the file handle
func OpenFileForStdout(filePath string, truncate bool) (*os.File, error) {
	// Ensure dir path exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	flag := os.O_WRONLY | os.O_CREATE
	if truncate {
		flag = flag | os.O_TRUNC
	} else {
		// if no truncate, we append (used for >> )
		flag = flag | os.O_APPEND
	}

	return os.OpenFile(filePath, flag, 0644)
}

// GetStdStream parses the redirection operator (e.g., `1>` or `2>`) and returns the corresponding stream type (stdout or stderr)
func GetStdStream(input string, pos int) (int, error) {
	if pos >= len(input) {
		return types.DefaultStdStream, nil
	}

	if input[pos] == ' ' {
		// > or >> are treated as 1> or 1>>
		return types.DefaultStdStream, nil
	}

	return strconv.Atoi(string(input[pos]))
}
