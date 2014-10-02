package pages

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/srhnsn/go-utils/database"
	emailUtil "github.com/srhnsn/go-utils/email"
	"github.com/srhnsn/go-utils/log"
	"github.com/srhnsn/go-utils/misc"
	"github.com/srhnsn/go-utils/webapps"

	"github.com/srhnsn/stuwomails/model"
)

var emailPattern *regexp.Regexp = regexp.MustCompile("^[^@]+@[^@]+$")

type indexPageSubmitData struct {
	mailingListId uint32
	firstName     string
	lastName      string
	room          string
	email         string
}

func IndexPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	webapps.SendResponse("index", r, w)
}

func IndexPageSubmitForm(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := webapps.GetTemplateData(r)
	T := webapps.GetFutureT(r)

	requestData := getIndexPageSubmitData(r)
	ok := validateIndexPageSubmitData(r, w, requestData)

	if !ok {
		return
	}

	err := sendEmail(requestData.email, requestData.mailingListId)

	if err != nil {
		log.Error.Printf("Subscription email failed: %s", err)
		sendServerError(r, w)
		return
	}

	ipAddress := misc.GetProxiedIpAddress(r)

	saveNewRequest(model.Request{
		MailingListId: requestData.mailingListId,
		FirstName:     requestData.firstName,
		LastName:      requestData.lastName,
		Room:          requestData.room,
		Email:         requestData.email,
		IpAddress:     ipAddress,
	})

	data["email"] = requestData.email

	data["title"] = T("page_index_email_success_title")
	data["message"] = T("page_index_email_success_message")
	webapps.FlashMessage(r, w)
}

func getIndexPageSubmitData(r *http.Request) indexPageSubmitData {
	mailingListIdInt, err := strconv.ParseUint(r.PostFormValue("mailing_list_id"), 10, 32)

	if err != nil {
		mailingListIdInt = 0
	}

	data := indexPageSubmitData{
		mailingListId: uint32(mailingListIdInt),
		firstName:     strings.TrimSpace(r.PostFormValue("first_name")),
		lastName:      strings.TrimSpace(r.PostFormValue("last_name")),
		room:          strings.TrimSpace(r.PostFormValue("room")),
		email:         strings.TrimSpace(r.PostFormValue("email")),
	}

	return data
}

func validateIndexPageSubmitData(r *http.Request, w http.ResponseWriter, requestData indexPageSubmitData) (ok bool) {
	data := webapps.GetTemplateData(r)
	T := webapps.GetFutureT(r)
	ok = false

	if requestData.mailingListId == 0 {
		data["title"] = T("page_index_invalid_mailing_list_title")
		data["message"] = T("page_index_invalid_mailing_list_message")
		webapps.FlashMessage(r, w)
		return
	}

	if requestData.firstName == "" || requestData.lastName == "" || requestData.room == "" || requestData.email == "" {
		data["title"] = T("page_index_empty_fields_title")
		data["message"] = T("page_index_empty_fields_message")
		webapps.FlashMessage(r, w)
		return
	}

	if !emailPattern.MatchString(requestData.email) {
		data["title"] = T("page_index_invalid_email_title")
		data["message"] = T("page_index_invalid_email_message")
		webapps.FlashMessage(r, w)
		return
	}

	return true
}

func getSubscribeEmail(mailingListId uint32) (string, error) {
	var subscribeEmail string

	err := database.DB.QueryRow(`
        SELECT
            subscribe_address
        FROM
            mailing_list
        WHERE
            id = ?
        LIMIT
            1
    `, mailingListId).Scan(&subscribeEmail)

	switch {
	case err == sql.ErrNoRows:
		log.Warning.Printf("getSubscribeEmail(%d) returned zero rows", mailingListId)
		return "", err
	case err != nil:
		log.Error.Printf("getSubscribeEmail(%d) failed: %s", mailingListId, err)
		return "", err
	}

	return subscribeEmail, nil
}

func saveNewRequest(request model.Request) {
	_, err := database.DB.NamedExec(`
        INSERT INTO
            request ( mailing_list_id , first_name , last_name , room , email , ip_address , creation_date )
        VALUES
            ( :mailing_list_id , :first_name , :last_name , :room , :email , :ip_address , NOW() )
    `, request)

	if err != nil {
		log.Error.Printf("saveNewRequest() failed: %s", err)
	}
}

func sendEmail(email string, mailingListId uint32) error {
	email = strings.Replace(email, "@", "=", -1)
	subscribeEmail, err := getSubscribeEmail(mailingListId)

	if err != nil {
		return err
	}

	to := fmt.Sprintf(subscribeEmail, email)

	errChan := emailUtil.SendEmail(emailUtil.Email{
		To:      to,
		Subject: "",
		Text:    "",
	})

	return <-errChan
}

func sendServerError(r *http.Request, w http.ResponseWriter) {
	data := webapps.GetTemplateData(r)
	T := webapps.GetFutureT(r)

	data["title"] = T("page_index_server_error_title")
	data["message"] = T("page_index_server_error_message")

	w.WriteHeader(http.StatusInternalServerError)
	webapps.FlashMessage(r, w)
}
