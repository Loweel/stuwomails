package model

import (
	"time"
)

type CheckLog struct {
	Id           uint32
	Input        string
	Password     string
	IpAddress    string
	CreationDate time.Time
}
