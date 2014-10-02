package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/srhnsn/go-utils/webapps"

	"github.com/srhnsn/stuwomails/pages"
)

type RoutesConfig struct {
	Asset    func(name string) ([]byte, error)
	AssetDir func(name string) ([]string, error)
}

const staticFilesPath = "www"

func GetRoutes(config RoutesConfig) http.Handler {
	routes := httprouter.New()

	routes.GET("/", pages.IndexPage)
	routes.POST("/", pages.IndexPageSubmitForm)

	routes.GET("/check/", pages.CheckPage)
	routes.GET("/check/:password/", pages.CheckPage)
	routes.POST("/check/", pages.CheckPageSubmitForm)

	routes.GET("/language/:language/", pages.LanguagePage)

	routes.GET("/static/:time/*filepath", webapps.GetFileServer(webapps.FileServerConfig{
		Asset:           config.Asset,
		AssetDir:        config.AssetDir,
		StaticFilesPath: staticFilesPath,
	}))

	return routes
}
