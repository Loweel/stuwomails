package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/srhnsn/go-utils/database"
	"github.com/srhnsn/go-utils/email"
	"github.com/srhnsn/go-utils/i18n"
	"github.com/srhnsn/go-utils/log"
	"github.com/srhnsn/go-utils/misc"
	"github.com/srhnsn/go-utils/webapps"

	appConfig "github.com/srhnsn/stuwomails/config"
	"github.com/srhnsn/stuwomails/config/routes"
)

const configFilename = "config.yaml"

func main() {
	var config appConfig.ConfigType
	misc.LoadConfig(configFilename, Asset, &config)
	initEmailRegexp(&config)

	appConfig.Config = config
	initMiddleware(config)

	addr := fmt.Sprintf(":%d", config.Server.Port)
	handler := webapps.RequestHandlerFunc

	log.Info.Printf("Trying to listen on port %d", config.Server.Port)
	err := http.ListenAndServe(addr, handler)

	if err != nil {
		log.Error.Fatalf("Cannot listen on port %d: %s", config.Server.Port, err)
	}
}

func initEmailRegexp(config *appConfig.ConfigType) {
	config.EmailBlacklistPatterns = make([]*regexp.Regexp, len(config.EmailBlacklist))

	for i, email := range config.EmailBlacklist {
		if email == "" {
			log.Error.Fatalf("Got empty blacklist email address")
		}

		config.EmailBlacklistPatterns[i] = regexp.MustCompile(email)
	}
}

func initMiddleware(config appConfig.ConfigType) {
	database.InitDatabase(config.Database)
	email.InitEmails(config.Mail)

	i18n.InitI18n(i18n.I18nConfig{
		Asset:           Asset,
		AssetDir:        AssetDir,
		DefaultLanguage: config.App.DefaultLanguage,
		Languages:       config.Languages,
	})

	webapps.InitCsp(webapps.CspConfig{config.Server.Csp})

	myroutes := routes.GetRoutes(routes.RoutesConfig{
		Asset:    Asset,
		AssetDir: AssetDir,
	})

	webapps.InitSessions(webapps.SessionConfig{
		Routes: myroutes,
		Secret: config.Server.Secret,
	})

	webapps.InitTemplates(webapps.TemplateConfig{
		Asset: Asset,
		Root:  config.Server.Root,
	})
}
