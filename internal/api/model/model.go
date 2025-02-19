package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Library struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string
}

type Users struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name          string
	Email         string
	ContactNumber string
	Role          string
	Library       *Library `gorm:"foreignKey:LibID"`
	LibID         *uuid.UUID
}

type BookInventory struct {
	ISBN            string   `gorm:"type:varchar(17);primaryKey"`
	Library         *Library `gorm:"foreignKey:LibID"`
	LibID           *uuid.UUID
	Title           string
	Authors         string
	Publisher       string
	Version         string
	TotalCopies     uint `gorm:"default:0"`
	AvailableCopies uint `gorm:"default:0"`
}

type RequestEvents struct {
	ReqID         uuid.UUID      `gorm:"type:uuid;primaryKey"`
	BookInventory *BookInventory `gorm:"foreignKey:BookID"`
	BookID        string
	Reader        *Users `gorm:"foreignKey:ReaderID"`
	ReaderID      uuid.UUID
	RequestDate   time.Time
	ApprovalDate  sql.NullTime
	Admin         *Users `gorm:"foreignKey:ApproverID"`
	ApproverID    *uuid.UUID
	RequestType   string
}

type IssueRegistry struct {
	IssueID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	BookInventory      *BookInventory `gorm:"foreignKey:ISBN"`
	ISBN               string
	Reader             *Users `gorm:"foreignKey:ReaderID"`
	ReaderID           *uuid.UUID
	AdminIssue         *Users `gorm:"foreignKey:IssueApproverID"`
	IssueApproverID    *uuid.UUID
	IssueStatus        string
	IssueDate          time.Time
	ExpectedReturnDate time.Time
	ReturnDate         sql.NullTime
	AdminReturn        *Users `gorm:"foreignKey:ReturnApproverID"`
	ReturnApproverID   *uuid.UUID
}
