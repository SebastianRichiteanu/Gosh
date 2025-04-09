package config

const (
	defaultPromptSymbol       = "$"
	defaultLogLevel           = "INFO"
	defaultEnableAutoComplete = true

	defaultGoshHomePath   = "~/.gosh"
	defaultLogFile        = "gosh.log"
	defaultHistoryFile    = "history"
	defaultMaxHistorySize = 10
	defaultGoshrcFile     = "goshrc"
)

const (
	envVarPromptSymbol       = "GOSH_SHELL_SYMBOL"
	envVarLogLevel           = "GOSH_LOG_LEVEL"
	envVarEnableAutoComplete = "GOSH_ENABLE_AUTOCOMPLETE"
	envVarLogFile            = "GOSH_LOG_FILE"
	envVarHistoryFile        = "GOSH_HISTORY_FILE"
	envVarMaxHistorySize     = "GOSH_MAX_HISTORY_SIZE"
	envVarGoshHomePath       = "GOSH_CONFIG_HOME"
)
