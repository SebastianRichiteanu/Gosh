package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/SebastianRichiteanu/Gosh/internal/utils"
)

type Config struct {
	PromptSymbol       string
	LogLevel           string
	LogFile            string
	GoshHomePath       string
	AliasFile          string
	HistoryFile        string
	MaxHistorySize     int
	EnableAutoComplete bool
	reloadCfgChannel   chan bool
}

func NewConfig(reloadCfgChannel chan bool) (*Config, error) {
	cfg := Config{
		PromptSymbol:       defaultPromptSymbol,
		LogLevel:           defaultLogLevel,
		LogFile:            defaultLogFile,
		HistoryFile:        defaultHistoryFile,
		MaxHistorySize:     defaultMaxHistorySize,
		EnableAutoComplete: defaultEnableAutoComplete,
		GoshHomePath:       defaultGoshHomePath,
		AliasFile:          defaultAliasFile,

		reloadCfgChannel: reloadCfgChannel,
	}

	goshrcFilePath := filepath.Join(cfg.GoshHomePath, defaultGoshrcFile)
	goshrcExpandedFilePath, err := utils.ExpandHomePath(goshrcFilePath)
	if err != nil {
		return nil, err
	}

	if err := cfg.EnsureConfig(goshrcExpandedFilePath); err != nil {
		return nil, fmt.Errorf("failed to ensure config exist: %v", err)
	}

	if err := utils.SourceFile(goshrcExpandedFilePath); err != nil {
		return nil, err
	}

	if err := cfg.Update(); err != nil {
		return nil, fmt.Errorf("failed to update config: %v", err)
	}

	go cfg.listenRefreshChan()

	return &cfg, nil
}

func (c *Config) Close() {
	close(c.reloadCfgChannel)
}

// EnsureConfig makes sure ~/.gosh/goshrc exists, creating it with defaults if needed.
func (c *Config) EnsureConfig(cfgFile string) error {
	cfgDir, err := utils.ExpandHomePath(c.GoshHomePath)
	if err != nil {
		return fmt.Errorf("could not expand home path: %v", err)
	}

	if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cfgDir, 0o755); err != nil {
			return fmt.Errorf("could not create config dir: %w", err)
		}
	}

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		if err := os.WriteFile(cfgFile, []byte(defaultConfig), 0644); err != nil {
			return fmt.Errorf("could not write default config: %w", err)
		}
	}

	return nil
}

func (c *Config) listenRefreshChan() {
	for {
		<-c.reloadCfgChannel
		err := c.Update()
		if err != nil {
			fmt.Printf("failed to refresh config: %v", err)
		}
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
	if envMaxHistorySize, exists := os.LookupEnv(envVarMaxHistorySize); exists {
		if maxHistorySize, err := strconv.Atoi(envMaxHistorySize); err == nil {
			c.MaxHistorySize = maxHistorySize
		} else {
			return fmt.Errorf("invalid value for MaxHistorySize: %v", err)
		}
	}

	if envAliasFile, exists := os.LookupEnv(envVarAliasFile); exists {
		c.AliasFile = envAliasFile
	}

	if envAutoComplete, exists := os.LookupEnv(envVarEnableAutoComplete); exists {
		c.EnableAutoComplete = envAutoComplete == "true"
	}

	if !filepath.IsAbs(c.LogFile) {
		c.LogFile = filepath.Join(c.GoshHomePath, c.LogFile)
	}
	if !filepath.IsAbs(c.HistoryFile) {
		c.HistoryFile = filepath.Join(c.GoshHomePath, c.HistoryFile)
	}
	if !filepath.IsAbs(c.AliasFile) {
		c.AliasFile = filepath.Join(c.GoshHomePath, c.AliasFile)
	}
	return nil
}
