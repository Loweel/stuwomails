package model

import (
	"time"
)

type Request struct {
	Id             uint32
	MailingListId  uint32
	FirstName      string
	LastName       string
	Room           string
	Email          string
	IpAddress      string
	ApprovalStatus string
	CreationDate   time.Time
}
