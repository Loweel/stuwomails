package pages

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/srhnsn/stuwomails/config"
)

func LanguagePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	language := ps.ByName("language")

	cookie := http.Cookie{
		Name:    "lang",
		Value:   language,
		Path:    config.Config.Server.Root,
		Expires: time.Now().Add(time.Hour * 24 * 365),
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, config.Config.Server.Root, 303)
}
