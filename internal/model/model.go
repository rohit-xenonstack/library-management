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
	Library       *Library `gorm:"foreignKey:LibID;association_foreignkey:ID"`
	LibID         *string
}

type BookInventory struct {
	ISBN            string `gorm:"type:varchar(17);primaryKey"`
	LibID           *string
	Title           string
	Authors         string
	Publisher       string
	Version         string
	TotalCopies     uint     `gorm:"default:0"`
	AvailableCopies uint     `gorm:"default:0"`
	Library         *Library `gorm:"foreignKey:LibID;association_foreignkey:ID"`
}

type RequestEvents struct {
	ReqID         string `gorm:"type:uuid;primaryKey"`
	BookID        string
	ReaderID      string
	RequestDate   time.Time
	ApprovalDate  sql.NullTime
	ApproverID    sql.NullString
	RequestType   string
	BookInventory *BookInventory `gorm:"foreignKey:BookID;association_foreignkey:ISBN"`
	Users         *Users         `gorm:"foreignKey:ReaderID,ApproverID;association_foreignkey:ID"`
}

type IssueRegistry struct {
	IssueID            string `gorm:"type:uuid;primaryKey"`
	ISBN               string
	ReaderID           string
	IssueApproverID    string
	IssueStatus        string
	IssueDate          time.Time
	ExpectedReturnDate time.Time
	ReturnDate         sql.NullTime
	ReturnApproverID   *string
	BookInventory      *BookInventory `gorm:"foreignKey:BookID;association_foreignkey:ISBN"`
	Users              *Users         `gorm:"foreignKey:ReaderID,IssueApproverID,ReturnApproverID;association_foreignkey:ID"`
}
