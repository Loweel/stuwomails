package pages

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/srhnsn/go-utils/database"
	"github.com/srhnsn/go-utils/log"
	"github.com/srhnsn/go-utils/misc"
	"github.com/srhnsn/go-utils/webapps"

	"github.com/srhnsn/stuwomails/config"
	"github.com/srhnsn/stuwomails/model"
)

type emailSearchResult struct {
	model.Request

	MailingListName  string
	CreationDateNice string
}

func CheckPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	data := webapps.GetTemplateData(r)

	data["password"] = ps.ByName("password")
	data["email"] = ""

	webapps.SendResponse("check", r, w)
}

func CheckPageSubmitForm(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := webapps.GetTemplateData(r)
	T := webapps.GetFutureT(r)
	StringT := webapps.GetStringT(r)

	password := r.PostFormValue("password")
	email := strings.TrimSpace(r.PostFormValue("email"))

	ipAddress := misc.GetProxiedIpAddress(r)

	logSearch(email, password, ipAddress)

	if password != config.Config.App.CheckPassword {
		log.Warning.Printf("Wrong check request password (%s) by %s", password, ipAddress)
		data["title"] = T("page_check_wrong_password_title")
		data["message"] = T("page_check_wrong_password_message")
		webapps.FlashMessage(r, w)
		return
	}

	requests := getEmailSearchResults(email, StringT("page_check_results_table_date_format"))

	data["count"] = len(requests)
	data["email"] = email
	data["password"] = password
	data["query_sent"] = true
	data["requests"] = requests

	webapps.SendResponse("check", r, w)
}

func getEmailSearchResults(email string, format string) []emailSearchResult {
	var result []emailSearchResult

	err := database.DB.Select(&result, `
        SELECT
            request.* , mailing_list.name AS mailing_list_name
        FROM
            request
        INNER JOIN
            mailing_list ON mailing_list.id = request.mailing_list_id
        WHERE
            request.email = ?
    `, email)

	if err != nil {
		log.Error.Printf("getEmailSearchResults(%s) failed: %s", email, err)
		return []emailSearchResult{}
	}

	for i, v := range result {
		result[i].CreationDateNice = v.CreationDate.Format(format)
	}

	return result
}

func logSearch(input, password, ipAddress string) {
	log.Trace.Printf("Got check request for %s by %s", input, ipAddress)

	_, err := database.DB.Exec(`
        INSERT INTO
            check_log ( input , password , ip_address , creation_date )
        VALUES
            ( ? , ? , ? , NOW() )
    `, input, password, ipAddress)

	if err != nil {
		log.Error.Printf("saveNewRequest() failed: %s", err)
	}
}
