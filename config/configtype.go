package config

import (
	"regexp"

	"github.com/srhnsn/go-utils/database"
	"github.com/srhnsn/go-utils/email"
)

type ConfigType struct {
	App struct {
		CheckPassword   string `mapstructure:"check_password"`
		DefaultLanguage string `mapstructure:"default_language"`
	}

	Database database.DatabaseConfig

	EmailBlacklist         []string `mapstructure:"email_blacklist"`
	EmailBlacklistPatterns []*regexp.Regexp

	Languages []string

	Mail email.EmailConfig

	Server struct {
		Csp    []string
		Debug  bool
		Port   uint16
		Root   string
		Secret string
	}
}

var Config ConfigType
