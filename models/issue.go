package models

import (
	"database/sql"
	"time"
)

type IssueRegistry struct {
	IssueID            string `gorm:"primaryKey"`
	ISBN               string
	ReaderID           string
	IssueApproverID    string
	IssueStatus        string
	IssueDate          time.Time
	ExpectedReturnDate time.Time
	ReturnDate         sql.NullTime
	ReturnApproverID   sql.NullString
}
