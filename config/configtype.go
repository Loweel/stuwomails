package config

import (
	"github.com/srhnsn/go-utils/database"
	"github.com/srhnsn/go-utils/email"
)

type ConfigType struct {
	App struct {
		CheckPassword   string `mapstructure:"check_password"`
		DefaultLanguage string `mapstructure:"default_language"`
	}

	Database database.DatabaseConfig

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
