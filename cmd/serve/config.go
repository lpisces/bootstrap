package serve

import (
	//"github.com/labstack/gommon/log"
	"gopkg.in/ini.v1"
)

var (
	Conf  *Config
	Debug bool
	Embed bool
)

type (
	// Config app config
	Config struct {
		Mode   string
		DB     *DBConfig
		Srv    *SrvConfig
		Site   *SiteConfig
		Secret *SecretConfig
		SMTP   *SMTPConfig
	}

	// SMTPConfig
	SMTPConfig struct {
		Hostname string
		Port     string
		Username string
		Password string
		FromAddr string
		FromName string
	}

	// DBConfig database config
	DBConfig struct {
		Driver     string
		DataSource string
	}

	// SrvConfig server config
	SrvConfig struct {
		Host string
		Port string
	}

	// SiteConfig site config
	SiteConfig struct {
		Name        string
		BaseURL     string
		SessionName string
	}

	SecretConfig struct {
		Session  string
		Password string
	}
)

func init() {
	Conf = DefaultConfig()
	Embed = false
}

// DefaultConfig get default config
func DefaultConfig() (config *Config) {

	db := &DBConfig{
		Driver:     "sqlite3",
		DataSource: "./bootstrap.db",
	}

	srv := &SrvConfig{
		Host: "0.0.0.0",
		Port: "1323",
	}

	site := &SiteConfig{
		Name:        "Bootstrap",
		BaseURL:     "http://127.0.0.1/",
		SessionName: "bs_sess",
	}

	secret := &SecretConfig{
		Session:  "secret",
		Password: "secret",
	}

	smtp := &SMTPConfig{
		Hostname: "",
		Port:     "",
		Username: "",
		Password: "",
		FromAddr: "",
		FromName: "",
	}

	config = &Config{
		Mode:   "development",
		DB:     db,
		Srv:    srv,
		Site:   site,
		Secret: secret,
		SMTP:   smtp,
	}
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
	db.DataSource = cfg.Section("db").Key("source").String()
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

	// site
	site := &SiteConfig{}
	site.Name = cfg.Section("site").Key("name").String()
	if site.Name == "" {
		site.Name = config.Site.Name
	}
	site.BaseURL = cfg.Section("site").Key("base_url").String()
	if site.BaseURL == "" {
		site.BaseURL = config.Site.BaseURL
	}
	site.SessionName = cfg.Section("site").Key("session_name").String()
	if site.SessionName == "" {
		site.SessionName = config.Site.SessionName
	}

	// secret
	secret := &SecretConfig{}
	secret.Session = cfg.Section("secret").Key("session").String()
	if secret.Session == "" {
		secret.Session = config.Secret.Session
	}
	secret.Password = cfg.Section("secret").Key("password").String()
	if secret.Password == "" {
		secret.Password = config.Secret.Password
	}

	// smtp
	smtp := &SMTPConfig{}
	smtp.Hostname = cfg.Section("smtp").Key("hostname").String()
	if smtp.Hostname == "" {
		smtp.Hostname = config.SMTP.Hostname
	}
	smtp.Port = cfg.Section("smtp").Key("port").String()
	if smtp.Port == "" {
		smtp.Port = config.SMTP.Port
	}
	smtp.Username = cfg.Section("smtp").Key("username").String()
	if smtp.Username == "" {
		smtp.Username = config.SMTP.Username
	}
	smtp.Password = cfg.Section("smtp").Key("password").String()
	if smtp.Password == "" {
		smtp.Password = config.SMTP.Password
	}
	smtp.FromAddr = cfg.Section("smtp").Key("from_addr").String()
	if smtp.FromAddr == "" {
		smtp.FromAddr = config.SMTP.FromAddr
	}
	smtp.FromName = cfg.Section("smtp").Key("from_name").String()
	if smtp.FromName == "" {
		smtp.FromName = config.SMTP.FromName
	}

	config.Mode = mode
	config.DB = db
	config.Srv = srv
	config.Site = site
	config.Secret = secret
	config.SMTP = smtp

	return
}
