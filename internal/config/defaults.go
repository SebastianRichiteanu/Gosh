package config

const (
	defaultPromptSymbol        = "$"
	defaultLogLevel            = "INFO"
	defaultEnabledAutoComplete = true

	defaultGoshHomePath = "~/.gosh"
	defaultLogFile      = "gosh.log"
	defaultHistoryFile  = "history"
	defaultGoshrcFile   = "goshrc"
)

const (
	envVarPromptSymbol        = "GOSH_SHELL_SYMBOL"
	envVarLogLevel            = "GOSH_LOG_LEVEL"
	envVarEnabledAutoComplete = "GOSH_AUTO_COMPLETE"
	envVarLogFile             = "GOSH_LOG_FILE"
	envVarHistoryFile         = "GOSH_HISTORY_FILE"
	envVarGoshHomePath        = "GOSH_CONFIG_HOME"
)
