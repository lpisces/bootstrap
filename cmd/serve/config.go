package serve

import (
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
		Mail   *MailConfig
	}

	// MailConfig
	MailConfig struct {
		Hostname string
		Port     string
		Username string
		Password string
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

	mail := &MailConfig{
		Hostname: "",
		Port:     "",
		Username: "",
		Password: "",
	}

	config = &Config{
		Mode:   "development",
		DB:     db,
		Srv:    srv,
		Site:   site,
		Secret: secret,
		Mail:   mail,
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

	// mail
	mail := &MailConfig{}
	mail.Hostname = cfg.Section("mail").Key("hostname").String()
	if mail.Hostname == "" {
		mail.Hostname = config.Mail.Hostname
	}
	mail.Port = cfg.Section("mail").Key("port").String()
	if mail.Port == "" {
		mail.Port = config.Mail.Port
	}
	mail.Username = cfg.Section("mail").Key("username").String()
	if mail.Username == "" {
		mail.Username = config.Mail.Username
	}
	mail.Password = cfg.Section("mail").Key("password").String()
	if mail.Password == "" {
		mail.Password = config.Mail.Password
	}

	config.Mode = mode
	config.DB = db
	config.Srv = srv
	config.Site = site
	config.Secret = secret
	config.Mail = mail

	return
}
