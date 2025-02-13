package models

import (
	"database/sql"
	"time"
)
type RequestEvents struct {
	ReqID        string `gorm:"primaryKey"`
	BookID       string
	ReaderID     string
	RequestDate  time.Time
	ApprovalDate sql.NullTime
	ApproverID   sql.NullString
	RequestType  string
}
