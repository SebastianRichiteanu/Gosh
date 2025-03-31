package tests

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/creack/pty"
	"github.com/stretchr/testify/assert"
)

const (
	goshTestBinaryPath = "./tmp/gosh"
)

type shellTest struct {
	name    string
	input   []string
	want    []string
	wantErr bool
}

func TestShell(t *testing.T) {
	pwdOutput, err := exec.Command("pwd").Output()
	if err != nil {
		t.Fatalf("Failed to execute pwd command to build tests: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get the user home dir to build tests: %v", err)
	}

	currentDir := strings.Trim(string(pwdOutput), "\r\n")

	tests := []shellTest{
		{
			name:    "test invalid command",
			input:   []string{"echo2 Hello, Gosh!"},
			want:    []string{"echo2: not found"},
			wantErr: false,
		},
		{
			name:    "test echo",
			input:   []string{"echo Hello, Gosh!"},
			want:    []string{"Hello, Gosh!"},
			wantErr: false,
		},
		{
			name:    "test type builtin",
			input:   []string{"type echo"},
			want:    []string{"echo is a shell builtin"},
			wantErr: false,
		},
		{
			name:    "test type executable",
			input:   []string{"type ls"},
			want:    []string{"ls is /usr/bin/ls"},
			wantErr: false,
		},
		{
			name:    "test type not found",
			input:   []string{"type ls2"},
			want:    []string{"ls2: not found"},
			wantErr: false,
		},
		{
			name:    "test exit",
			input:   []string{"exit"},
			want:    []string{},
			wantErr: true,
		},
		{
			name:    "test pwd",
			input:   []string{"pwd"},
			want:    []string{currentDir},
			wantErr: false,
		},
		{
			name:    "test cd forward and backwords",
			input:   []string{"cd ./tmp", "pwd", "cd ..", "pwd"},
			want:    []string{"", currentDir + "/tmp", "", currentDir},
			wantErr: false,
		},
		{
			name:    "test cd random dir",
			input:   []string{"cd /tmp", "pwd"},
			want:    []string{"", "/tmp"},
			wantErr: false,
		},
		{
			name:    "test cd home dir",
			input:   []string{"cd ~", "pwd"},
			want:    []string{"", homeDir},
			wantErr: false,
		},
		{
			name:    "test run program - execute shell in shell",
			input:   []string{"./tmp/gosh", "echo 123"},
			want:    []string{"", "123"},
			wantErr: false,
		},
		{
			name:    "test single quotes",
			input:   []string{"echo 'Hello,    Gosh!'"},
			want:    []string{"Hello,    Gosh!"},
			wantErr: false,
		},
		{
			name:    "test double quotes",
			input:   []string{`echo "Hello,  Gosh!"  "Hi"`},
			want:    []string{"Hello,  Gosh! Hi"},
			wantErr: false,
		},
		{
			name:    "test double quotes",
			input:   []string{`echo "Hello,  Go'sh!"  "Hi"`},
			want:    []string{"Hello,  Go'sh! Hi"},
			wantErr: false,
		},
		{
			name:    "test backslash outside quotes",
			input:   []string{`echo "Hello\  Gosh!"`},
			want:    []string{`Hello\  Gosh!`},
			wantErr: false,
		},
		{
			name:    "test backslash outside quotes 2",
			input:   []string{`echo Hello\ \ Gosh!`},
			want:    []string{"Hello  Gosh!"},
			wantErr: false,
		},
		{
			name:    "test backslash within single quotes",
			input:   []string{`echo 'Hello\ Gosh!'`},
			want:    []string{`Hello\ Gosh!`},
			wantErr: false,
		},
		{
			name:    "test backslash within double quotes",
			input:   []string{`echo "Hello\ 'Gosh'!"`},
			want:    []string{`Hello\ 'Gosh'!`},
			wantErr: false,
		},
		{
			name:    "test backslash within double quotes 2",
			input:   []string{`echo "Hello \"Gosh\"!"`},
			want:    []string{`Hello "Gosh"!`},
			wantErr: false,
		},
		{
			name: "test executing a quoted executable",
			input: []string{
				`go build -o './tmp/name with "quotes"' ../cmd/gosh/main.go`, // build shell with quotes
				`./tmp/'name with "quotes"'`,                                 // execute and enter shell
				"pwd",                                                        // test new shell
				"exit",                                                       // exit second shell, this shouldn't be an error because we fall back to the first shell
			},
			want:    []string{"", "", currentDir, ""},
			wantErr: false,
		},
		{
			name:    "test stdout redirection",
			input:   []string{"echo Hello, Gosh! > ./tmp/output.txt", "cat ./tmp/output.txt"},
			want:    []string{"", "Hello, Gosh!"},
			wantErr: false,
		},
		{
			name:    "test stderr redirection",
			input:   []string{"ls non_existing_file 2> ./tmp/error.txt", "cat ./tmp/error.txt"},
			want:    []string{"", "ls: cannot access 'non_existing_file': No such file or directory"},
			wantErr: false,
		},
		{
			name:    "test append stdout",
			input:   []string{"echo First Line > ./tmp/output2.txt", "echo Second Line >> ./tmp/output2.txt", "cat ./tmp/output2.txt"},
			want:    []string{"", "", "First Line\r\nSecond Line"},
			wantErr: false,
		},
		{
			name:    "test append stderr",
			input:   []string{"ls non_existing_file 2> ./tmp/error2.txt", "ls another_missing_file 2>> ./tmp/error2.txt", "cat ./tmp/error2.txt"},
			want:    []string{"", "", "ls: cannot access 'non_existing_file': No such file or directory\r\nls: cannot access 'another_missing_file': No such file or directory"},
			wantErr: false,
		},

		// TODO: create tests for autocompletion
		// TODO: create tests for history (up/down arrow)
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if !test.wantErr {
				assert.Equal(t, len(test.input), len(test.want), "the number of input cmds does not match with the number of outputs")
			}

			ptyMaster, err := initShell()
			if err != nil {
				t.Fatal(err)
			}
			defer ptyMaster.Close()

			for idx, inputCmd := range test.input {
				if err := sendCommand(ptyMaster, inputCmd); err != nil {
					t.Errorf("Failed to write to PTY: %v", err)
				}

				time.Sleep(100 * time.Millisecond)

				got, err := readOutput(ptyMaster)
				if (err != nil) != test.wantErr {
					t.Errorf("wrong output error=%v, wantErr %v", err, test.wantErr)
					return
				}

				if !test.wantErr {
					assert.Equal(t, test.want[idx], got, "output does not match")
				}

			}
		})
	}
}

// initShell will init a pty shell for gosh
func initShell() (*os.File, error) {
	cmd := exec.Command(goshTestBinaryPath)
	ptyMaster, err := pty.Start(cmd)
	if err != nil {
		return nil, fmt.Errorf("Failed to start shell in PTY: %w", err)
	}

	time.Sleep(500 * time.Millisecond)

	return ptyMaster, nil
}

// sendCommand sends a command to the PTY
func sendCommand(ptyMaster *os.File, command string) error {
	_, err := ptyMaster.Write([]byte(command + "\r\n"))
	return err
}

// readOutput reads the output from the PTY until the prompt is reached
func readOutput(ptyMaster *os.File) (string, error) {
	reader := bufio.NewReader(ptyMaster)

	// Read and ignore the prompt line
	_, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Read until the next prompt
	output, err := reader.ReadString('$') // TODO: add config for tests
	if err != nil {
		return "", err
	}

	outputStr := strings.Trim(output, "\r\n$")
	return outputStr, nil
}
