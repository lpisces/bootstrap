package serve

import (
	"gopkg.in/ini.v1"
)

var (
	Conf  *Config
	Debug bool
)

// Config app config
type Config struct {
	Mode string
	DB   *DBConfig
	Srv  *SrvConfig
}

// DBConfig database config
type DBConfig struct {
	Driver     string
	DataSource string
}

// SrvConfig server config
type SrvConfig struct {
	Host string
	Port string
}

// DefaultConfig get default config
func DefaultConfig() (config *Config) {
	config = &Config{}
	db := &DBConfig{}
	db.Driver = "sqlite3"
	db.DataSource = "./bootstrap.db"

	srv := &SrvConfig{}
	srv.Host = "127.0.0.1"
	srv.Port = "1323"

	config.Mode = "development"
	config.DB = db
	config.Srv = srv
	return
}

// Load load config from file override default config
func (config *Config) Load(path string) (err error) {

	// load from ini file
	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}

	// run mode
	mode := cfg.Section("").Key("mode").In("development", []string{"development", "production", "testing"})

	// database
	db := &DBConfig{}
	db.Driver = cfg.Section("db").Key("driver").In("sqlite3", []string{"sqlite3", "mysql"})
	db.DataSource = cfg.Section("db").Key("dataSource").String()
	if db.DataSource == "" {
		db.DataSource = config.DB.DataSource
	}

	// server
	srv := &SrvConfig{}
	srv.Host = cfg.Section("srv").Key("host").String()
	if srv.Host == "" {
		srv.Host = config.Srv.Host
	}
	srv.Port = cfg.Section("srv").Key("port").String()
	if srv.Port == "" {
		srv.Port = config.Srv.Port
	}

	config.Mode = mode
	config.DB = db
	config.Srv = srv

	return
}
