package pages

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/srhnsn/go-utils/database"
	"github.com/srhnsn/go-utils/i18n"
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

type checkPageSubmitData struct {
	password       string
	email          string
	requestId      uint32
	approvalStatus string
}

const approvalStatusLanguagePrefix = "page_check_approval_status_"

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

	requestData := getCheckPageSubmitData(r)
	ipAddress := misc.GetProxiedIpAddress(r)

	logSearch(requestData.email, requestData.password, ipAddress)

	if requestData.password != config.Config.App.CheckPassword {
		log.Warning.Printf("Wrong check request password (%s) by %s", requestData.password, ipAddress)
		data["title"] = T("page_check_wrong_password_title")
		data["message"] = T("page_check_wrong_password_message")
		webapps.FlashMessage(r, w)
		return
	}

	if requestData.approvalStatus != "" {
		saveNewApprovalStatus(requestData.requestId, requestData.approvalStatus)
	}

	requests := getEmailSearchResults(requestData.email, StringT("page_check_results_table_date_format"))
	requestsCount := len(requests)

	setCheckCountMessages(data, requestsCount, requestData, T)

	data["count"] = requestsCount
	data["email"] = requestData.email
	data["password"] = requestData.password
	data["query_sent"] = true
	data["requests"] = requests
	data["approval_status_language_prefix"] = approvalStatusLanguagePrefix

	webapps.SendResponse("check", r, w)
}

func getCheckPageSubmitData(r *http.Request) checkPageSubmitData {
	requestId, err := strconv.ParseUint(r.PostFormValue("id"), 10, 32)

	if err != nil {
		requestId = 0
	}

	data := checkPageSubmitData{
		requestId:      uint32(requestId),
		password:       r.PostFormValue("password"),
		email:          strings.TrimSpace(r.PostFormValue("email")),
		approvalStatus: r.PostFormValue("approval_status"),
	}

	return data
}

func getEmailSearchResults(email string, format string) []emailSearchResult {
	var err error
	var result []emailSearchResult
	var sql string

	sql = `
        SELECT
            request.* , mailing_list.name AS mailing_list_name
        FROM
            request
        INNER JOIN
            mailing_list ON mailing_list.id = request.mailing_list_id
        WHERE
            %s`

	if email == "" {
		sql = fmt.Sprintf(sql, `request.approval_status = "open"`)
		err = database.DB.Select(&result, sql)
	} else {
		sql = fmt.Sprintf(sql, "request.email = ?")
		err = database.DB.Select(&result, sql, email)
	}

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
	var what string

	if input == "" {
		what = "all open requests"
	} else {
		what = `"` + input + `"`
	}

	log.Trace.Printf(`Got check request for %s by %s`, what, ipAddress)

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

func saveNewApprovalStatus(id uint32, status string) {
	log.Trace.Printf("Setting approval status of request ID %d to %s", id, status)

	_, err := database.DB.Exec(`
        UPDATE
            request
        SET
            approval_status = ?
        WHERE
            id = ?
        LIMIT
            1
    `, status, id)

	if err != nil {
		log.Error.Printf("saveNewApprovalStatus() failed: %s", err)
	}
}

func setCheckCountMessages(data webapps.TemplateData, requestsCount int, requestData checkPageSubmitData, T i18n.FutureTranslateFunc) {
	if requestsCount == 0 {
		if requestData.email == "" {
			data["no_results_message"] = T("page_check_results_none_open")
			data["results_count_message"] = T("page_check_results_for_open")
		} else {
			data["no_results_message"] = T("page_check_results_none")
			data["results_count_message"] = T("page_check_results_for")
		}
	} else {
		if requestData.email == "" {
			data["results_count_message"] = T("page_check_results_for_open", requestsCount)
		} else {
			data["results_count_message"] = T("page_check_results_for", requestsCount)
		}
	}
}
