package serve

import (
	"gopkg.in/ini.v1"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
)

type Config struct {
	Mode     string
	Database *DatabaseConfig
	Server   *ServerConfig
}

type DatabaseConfig struct {
	Database string // database typo
	Dialect  string // database connection
}

type ServerConfig struct {
	Port     string // listen port
	Host     string // serve host
}

// loadConfig load config or create config file if not exist
func (conf *Config) Load(c *cli.Context) (err error) {
	path := c.String("config")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = "./config.default.ini"
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := ioutil.WriteFile(path, []byte(""), 0644)
		if err != nil {
			return err
		}
	}

	cfg, err := ini.Load(path)
	if err != nil {
		return
	}

	// run mode
	conf.Mode = cfg.Section("").Key("mode").In("development", []string{"development", "production", "testing"})

	// database
	conf.Database.Database = cfg.Section("database").Key("database").In("sqlite3", []string{"sqlite3", "mysql"})
	conf.Database.Dialect = cfg.Section("database").Key("dialect").String()

	// server
	conf.Server.Host = cfg.Section("server").Key("host").String()
	conf.Server.Port = cfg.Section("server").Key("port").String()

	return
}
