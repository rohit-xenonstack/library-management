package model

import (
	"database/sql"
	"time"
)

type Library struct {
	ID   string `gorm:"type:uuid;primaryKey"`
	Name string
}

type Users struct {
	ID            string `gorm:"type:uuid;primaryKey"`
	Name          string
	Email         string
	ContactNumber string
	Role          string
	Library       *Library `gorm:"foreignKey:LibID"`
	LibID         *string
}

type BookInventory struct {
	ISBN            string   `gorm:"type:varchar(17);primaryKey"`
	Library         *Library `gorm:"foreignKey:LibID"`
	LibID           *string
	Title           string
	Authors         string
	Publisher       string
	Version         string
	TotalCopies     uint `gorm:"default:0"`
	AvailableCopies uint `gorm:"default:0"`
}

type RequestEvents struct {
	ReqID         string         `gorm:"type:uuid;primaryKey"`
	BookInventory *BookInventory `gorm:"foreignKey:BookID"`
	BookID        string
	Reader        *Users `gorm:"foreignKey:ReaderID"`
	ReaderID      string
	RequestDate   time.Time
	ApprovalDate  sql.NullTime
	Admin         *Users `gorm:"foreignKey:ApproverID"`
	ApproverID    sql.NullString
	RequestType   string
}

type IssueRegistry struct {
	IssueID            string         `gorm:"type:uuid;primaryKey"`
	BookInventory      *BookInventory `gorm:"foreignKey:ISBN"`
	ISBN               string
	Reader             *Users `gorm:"foreignKey:ReaderID"`
	ReaderID           string
	AdminIssue         *Users `gorm:"foreignKey:IssueApproverID"`
	IssueApproverID    string
	IssueStatus        string
	IssueDate          time.Time
	ExpectedReturnDate time.Time
	ReturnDate         sql.NullTime
	AdminReturn        *Users `gorm:"foreignKey:ReturnApproverID"`
	ReturnApproverID   *string
}
