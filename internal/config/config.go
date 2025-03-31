package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SebastianRichiteanu/Gosh/internal/utils"
)

type Config struct {
	PromptSymbol       string
	LogLevel           string
	LogFile            string
	HistoryFile        string
	EnableAutoComplete bool
	GoshHomePath       string

	reloadCfgChannel chan bool
}

func NewConfig(reloadCfgChannel chan bool) (*Config, error) {
	cfg := Config{
		PromptSymbol:       defaultPromptSymbol,
		LogLevel:           defaultLogLevel,
		LogFile:            defaultLogFile,
		HistoryFile:        defaultHistoryFile,
		EnableAutoComplete: defaultEnableAutoComplete,
		GoshHomePath:       defaultGoshHomePath,

		reloadCfgChannel: reloadCfgChannel,
	}

	goshrcFilePath := filepath.Join(cfg.GoshHomePath, defaultGoshrcFile)
	goshrcExpandedFilePath, err := utils.ExpandHomePath(goshrcFilePath)
	if err != nil {
		return nil, err
	}

	if err := utils.SourceFile(goshrcExpandedFilePath); err != nil {
		return nil, err
	}

	cfg.Update()

	go cfg.listenRefreshChan()

	return &cfg, nil
}

func (c *Config) listenRefreshChan() {
	for {
		_ = <-c.reloadCfgChannel
		_ = c.Update() // TODO: treat error?
	}
}

func (c *Config) Update() error {
	if envGoshHomePath, exists := os.LookupEnv(envVarGoshHomePath); exists {
		c.GoshHomePath = envGoshHomePath
	}

	var err error
	if c.GoshHomePath, err = utils.ExpandHomePath(c.GoshHomePath); err != nil {
		return fmt.Errorf("Failed to expand home path: %v", err)
	}

	if err = os.MkdirAll(c.GoshHomePath, 0755); err != nil {
		return fmt.Errorf("Failed to create dir: %v", err)
	}

	if prompt, exists := os.LookupEnv(envVarPromptSymbol); exists {
		c.PromptSymbol = prompt
	}
	if envLogLevel, exists := os.LookupEnv(envVarLogLevel); exists {
		c.LogLevel = envLogLevel
	}
	if envLogFile, exists := os.LookupEnv(envVarLogFile); exists {
		c.LogFile = envLogFile
	}
	if envHistoryFile, exists := os.LookupEnv(envVarHistoryFile); exists {
		c.HistoryFile = envHistoryFile
	}
	if envAutoComplete, exists := os.LookupEnv(envVarEnableAutoComplete); exists {
		c.EnableAutoComplete = envAutoComplete == "true"
	}

	c.LogFile = filepath.Join(c.GoshHomePath, c.LogFile)
	c.HistoryFile = filepath.Join(c.GoshHomePath, c.HistoryFile)

	return nil
}
