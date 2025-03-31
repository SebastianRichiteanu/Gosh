package config

import (
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

func NewConfig(reloadCfgChannel chan bool) *Config {
	cfg := Config{
		PromptSymbol:       defaultPromptSymbol,
		LogLevel:           defaultLogLevel,
		LogFile:            defaultLogFile,
		HistoryFile:        defaultHistoryFile,
		EnableAutoComplete: defaultEnabledAutoComplete,
		GoshHomePath:       defaultGoshHomePath,

		reloadCfgChannel: reloadCfgChannel,
	}

	// TODO: on create source them
	// for _, file := range []string{} {
	// 	fp := filepath.Join(c.GoshHomePath, file)
	// } TODO: source rc, env

	cfg.Update()

	go cfg.listenRefreshChan()

	return &cfg
}

func (c *Config) listenRefreshChan() {
	for {
		_ = <-c.reloadCfgChannel
		c.Update()
	}
}

func (c *Config) Update() {
	if envGoshHomePath, exists := os.LookupEnv(envVarGoshHomePath); exists {
		c.GoshHomePath = envGoshHomePath
	}

	// TODO: treat error
	c.GoshHomePath, _ = utils.ExpandHomePath(c.GoshHomePath)

	err := os.MkdirAll(c.GoshHomePath, 0755)
	if err != nil {
		return // TODO: handle err
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
	if envAutoComplete, exists := os.LookupEnv(envVarEnabledAutoComplete); exists {
		c.EnableAutoComplete = envAutoComplete == "true"
	}

	c.LogFile = filepath.Join(c.GoshHomePath, c.LogFile)
	c.HistoryFile = filepath.Join(c.GoshHomePath, c.HistoryFile)
}
