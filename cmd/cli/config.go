package cli

import (
	"gopkg.in/ini.v1"
)

var (
	Conf *CliConfig
)

type (
	// CliConfig app config
	CliConfig struct {
		Env      string
		Debug    bool
		CacheDir string
		Source   string
		Output   string
		*DBConfig
		*Kuaidi100Config
	}

	// DBConfig database config
	DBConfig struct {
		Driver     string
		DataSource string
	}

	// kuaidi100
	Kuaidi100Config struct {
		Customer string
		Key      string
	}
)

func init() {
	Conf = (&CliConfig{}).Default()
}

// DefaultConfig get default config
func (config *CliConfig) Default() *CliConfig {

	config = &CliConfig{
		"dev",
		true,
		"./cache",
		"./package.csv",
		"./output.csv",
		&DBConfig{
			Driver:     "sqlite3",
			DataSource: "./data.db",
		},
		&Kuaidi100Config{
			Customer: "",
			Key:      "",
		},
	}

	return config
}

// Load load config from file override default config
func (config *CliConfig) LoadFromIni(path string) (err error) {
	return ini.MapTo(config, path)
}
