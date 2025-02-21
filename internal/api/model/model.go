package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Library struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string    `gorm:"unique"`
}

type Users struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"user_id" binding:"required"`
	Name          string     `gorm:"" json:"name" binding:"required"`
	Email         string     `gorm:"unique" json:"email" binding:"required"`
	ContactNumber string     `gorm:"" json:"contact" binding:"required"`
	Role          string     `gorm:"" json:"role" binding:"required"`
	Library       *Library   `gorm:"foreignKey:LibID;references:ID" json:"-"`
	LibID         *uuid.UUID `gorm:"" json:"library_id"`
}

type BookInventory struct {
	ISBN            string     `gorm:"type:varchar(20);primaryKey" json:"isbn" binding:"required"`
	Library         *Library   `gorm:"foreignKey:LibID;references:ID" json:"-"`
	LibID           *uuid.UUID `gorm:"" json:"library_id" binding:"required"`
	Title           string     `gorm:"" json:"title" binding:"required"`
	Authors         string     `gorm:"" json:"authors" binding:"required"`
	Publisher       string     `gorm:"" json:"publisher" binding:"required"`
	Version         string     `gorm:"" json:"version" binding:"required"`
	TotalCopies     uint       `gorm:"" json:"total_copies" binding:"required"`
	AvailableCopies uint       `gorm:"" json:"available_copies" binding:"required"`
}

type RequestEvents struct {
	ReqID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"request_id"`
	BookInventory *BookInventory `gorm:"foreignKey:BookID;references:ISBN" json:"-"`
	BookID        string         `gorm:"" json:"isbn"`
	Reader        *Users         `gorm:"foreignKey:ReaderID;references:ID" json:"-"`
	ReaderID      uuid.UUID      `gorm:"" json:"reader_id"`
	RequestDate   time.Time      `gorm:"" json:"request_date"`
	ApprovalDate  sql.NullTime   `gorm:"" json:"approval_date"`
	Admin         *Users         `gorm:"foreignKey:ApproverID;references:ID" json:"-"`
	ApproverID    *uuid.UUID     `gorm:"" json:"approver_id"`
	RequestType   string         `gorm:"" json:"request_type"`
}

type IssueRegistry struct {
	IssueID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	BookInventory      *BookInventory `gorm:"foreignKey:BookID;references:ISBN"`
	BookID             string         `gorm:""`
	Reader             *Users         `gorm:"foreignKey:ReaderID;references:ID"`
	ReaderID           *uuid.UUID     `gorm:""`
	AdminIssue         *Users         `gorm:"foreignKey:IssueApproverID;references:ID"`
	IssueApproverID    *uuid.UUID     `gorm:""`
	IssueStatus        string         `gorm:""`
	IssueDate          time.Time      `gorm:""`
	ExpectedReturnDate time.Time      `gorm:""`
	ReturnDate         sql.NullTime   `gorm:""`
	AdminReturn        *Users         `gorm:"foreignKey:ReturnApproverID;references:ID"`
	ReturnApproverID   *uuid.UUID     `gorm:""`
}
