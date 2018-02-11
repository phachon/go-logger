package go_logger

// logger config
type Config struct {
	Console *ConsoleConfig
	File *FileConfig
	Api *ApiConfig
}

func NewConfigConsole(consoleConfig *ConsoleConfig) *Config {
	return &Config{
		Console: consoleConfig,
	}
}

func NewConfigFile(fileConfig *FileConfig) *Config {
	return &Config{
		File: fileConfig,
	}
}

func NewConfigApi(apiConfig *ApiConfig) *Config {
	return &Config{
		Api: apiConfig,
	}
}