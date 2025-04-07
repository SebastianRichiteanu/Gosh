package closer

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/SebastianRichiteanu/Gosh/internal/config"
	"github.com/SebastianRichiteanu/Gosh/internal/logger"
	"github.com/SebastianRichiteanu/Gosh/internal/prompt"
)

type Closer struct {
	exitChannel   chan int
	osSignalsChan chan os.Signal

	prompt *prompt.Prompt
	cfg    *config.Config
	logger *logger.Logger
}

// NewCloser creates and returns a new Closer instance
func NewCloser(exitChannel chan int, prompt *prompt.Prompt, cfg *config.Config, logger *logger.Logger) *Closer {
	return &Closer{
		exitChannel:   exitChannel,
		osSignalsChan: make(chan os.Signal, 1),

		prompt: prompt,
		cfg:    cfg,
		logger: logger,
	}
}

// Recover catches panics, logs the error and stack trace, and gracefully exits the program
func (c *Closer) Recover() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "caught panic and recovered for gracefully exit: %v\n", r)
		debug.PrintStack()
		c.HandleExit(1)
	}
}

// ListenForSignals continuously listens for OS signals and internal exit codes to initiate a graceful shutdown
func (c *Closer) ListenForSignals() {
	signal.Notify(c.osSignalsChan,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)

	for {
		select {
		case sig := <-c.osSignalsChan:
			c.HandleSignal(sig)
		case code := <-c.exitChannel:
			c.HandleExit(code)
		}
	}
}

// HandleSignal maps an incoming OS signal to an appropriate exit code and triggers a shutdown
func (c *Closer) HandleSignal(sig os.Signal) {
	// By convention exit codes related to signals are often: code = 128 + signal number
	c.HandleExit(128 + int(sig.(syscall.Signal)))
}

// HandleExit performs cleanup of all components and terminates the program with the specified exit code
func (c *Closer) HandleExit(code int) {
	close(c.osSignalsChan)
	close(c.exitChannel)

	c.prompt.Close()
	c.cfg.Close()
	c.logger.Close()

	os.Exit(code)
}
