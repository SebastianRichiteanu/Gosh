package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/SebastianRichiteanu/Gosh/internal/types"
)

// FindPath searches for the given command in the system's PATH, or checks if it's a local or full path
// It returns the full file path if found or an empty string if not
func FindPath(cmd string) string {
	if len(cmd) == 0 {
		return ""
	}

	switch cmd[0] {
	case '.': // Local path, relative to curent dir
		cmd = cmd[1:]
		currentDir, err := os.Getwd()
		if err != nil {
			return "" // TODO: treat error? maybe once I add debug
		}

		fp := filepath.Join(currentDir, cmd)
		if _, err := os.Stat(fp); err == nil {
			return fp
		}
	case '/': // Absolute path
		if _, err := os.Stat(cmd); err == nil {
			return cmd
		}
	default: // Search in PATH
		paths := strings.Split(os.Getenv(types.PathEnvVar), types.PathDelimiter)
		for _, path := range paths {
			fp := filepath.Join(path, cmd)
			if _, err := os.Stat(fp); err == nil {
				return fp
			}
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

		// Check if original path ends with a slash
		hasTrailingSlash := strings.HasSuffix(path, string(os.PathSeparator))

		// Join path segments
		expandedPath := filepath.Join(homeDir, path[1:])

		// Re-add trailing slash if it was present
		if hasTrailingSlash && !strings.HasSuffix(expandedPath, string(os.PathSeparator)) {
			expandedPath += string(os.PathSeparator)
		}

		return expandedPath, nil
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

// BlockCtrlC will start a channel and will listen for OS signals
// and will ignore Ctrl+C to handle exit gracefully
func BlockCtrlC() {
	// Create a channel to listen for incoming OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT) // Listen for SIGINT (Ctrl+C)

	go func() {
		<-sigChan
		log.Println("Ctrl+C caught, but not exiting...")
	}()
}

// FindLongestPrefix takes a list of strings and returns the longest common prefix among them
func FindLongestPrefix(cmds []string) string {
	common := ""
	isCommon := true

	for i := 0; isCommon; i++ {
		var current byte
		for j := 0; j < len(cmds); j++ {
			if i >= len(cmds[j]) {
				isCommon = false
				break
			} else if j == 0 {
				current = cmds[j][i]
			} else if current != cmds[j][i] {
				isCommon = false
				break
			}
		}
		if isCommon {
			common += string(current)
		}
	}
	return common
}

func HandleExportLine(line string) {
	line = strings.TrimPrefix(line, "export")
	line = strings.Trim(line, " ")

	parts := strings.SplitN(line, "=", 2)
	key := strings.TrimSpace(parts[0])
	value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)

	os.Setenv(key, value)
}

func SourceFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Failed to open %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check if it's a variable assignment (e.g., VAR=value)
		if strings.Contains(line, "=") {
			HandleExportLine(line)
			continue
		}

		// Execute command if not
		cmd := exec.Command("sh", "-c", line)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error executing command:", line, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading %s: %w", filePath, err)
	}

	return nil
}

func ExitShell(exitCode int) {
	//logger.Close()
	// Maybe create chan that waits for signals for stuff that I want to execute before exit

	os.Exit(exitCode)
}
