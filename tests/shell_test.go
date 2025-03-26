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

func TestShell(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    []string
		wantErr bool
	}{
		// TODO: I think it's best to move these somewhere else

		{name: "test invalid command", input: []string{"echo2 Hello, Gosh!"}, want: []string{"echo2: not found"}},
		{name: "test echo", input: []string{"echo Hello, Gosh!"}, want: []string{"Hello, Gosh!"}},
		{name: "test type builtin", input: []string{"type echo"}, want: []string{"echo is a shell builtin"}},
		{name: "test type executable", input: []string{"type ls"}, want: []string{"ls is /usr/bin/ls"}},
		{name: "test exit", input: []string{"exit"}, wantErr: true},

		{name: "test pwd", input: []string{"pwd"}, want: []string{"/mnt/d/Programming/Programming/GitHub/Gosh/tests"}}, // TODO: create  a main shell where these get executed and compare results :)

		{name: "test cd", input: []string{"cd ..", "pwd"}, want: []string{"", "/mnt/d/Programming/Programming/GitHub/Gosh"}}, // TODO:: add more

		{name: "test run program", input: []string{"ls ./tmp"}, want: []string{"gosh"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ptyMaster, err := initShell()
			if err != nil {
				t.Fatal(err)
			}
			defer ptyMaster.Close()

			for idx, inputCmd := range test.input {
				if err := sendCommand(ptyMaster, inputCmd); err != nil {
					t.Errorf("Failed to write to PTY: %v", err)
				}

				got, err := readOutput(ptyMaster)
				if (err != nil) != test.wantErr {
					t.Errorf("wrong output error=%v, wantErr %v", err, test.wantErr)
					return
				}

				if !test.wantErr {
					assert.Equal(t, got, test.want[idx], "test_name", test.name)
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
	output, err := reader.ReadString('$')
	if err != nil {
		return "", err
	}

	outputStr := strings.Trim(output, "\r\n$")
	return outputStr, nil
}
