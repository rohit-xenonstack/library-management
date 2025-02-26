package model

type Library struct {
	ID   string `gorm:"primaryKey" json:"library_id" binding:"required"`
	Name string `gorm:"unique" json:"name" binding:"required"`
}

type Users struct {
	ID            string   `gorm:"primaryKey" json:"user_id" binding:"required"`
	Name          string   `gorm:"" json:"name" binding:"required"`
	Email         string   `gorm:"unique" json:"email" binding:"required"`
	ContactNumber string   `gorm:"" json:"contact" binding:"required"`
	Role          string   `gorm:"" json:"role" binding:"required"`
	Library       *Library `gorm:"foreignKey:LibID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	LibID         *string  `gorm:"" json:"library_id"`
}

type BookInventory struct {
	ISBN            string   `gorm:"type:varchar(20);primaryKey" json:"isbn" binding:"required"`
	Library         *Library `gorm:"foreignKey:LibID;references:ID" json:"-"`
	LibID           *string  `gorm:"" json:"library_id" binding:"required"`
	Title           string   `gorm:"" json:"title" binding:"required"`
	Authors         string   `gorm:"" json:"authors" binding:"required"`
	Publisher       string   `gorm:"" json:"publisher" binding:"required"`
	Version         string   `gorm:"" json:"version" binding:"required"`
	TotalCopies     uint     `gorm:"" json:"total_copies" binding:"required"`
	AvailableCopies uint     `gorm:"" json:"available_copies" binding:"required"`
}

type RequestEvents struct {
	ReqID         string         `gorm:"primaryKey" json:"request_id" binding:"required"`
	BookInventory *BookInventory `gorm:"foreignKey:BookID;references:ISBN" json:"-"`
	BookID        string         `gorm:"" json:"isbn" binding:"required"`
	Reader        *Users         `gorm:"foreignKey:ReaderID;references:ID" json:"-"`
	ReaderID      string         `gorm:"" json:"reader_id" binding:"required"`
	RequestDate   string         `gorm:"" json:"request_date" binding:"required"`
	ApprovalDate  *string        `gorm:"" json:"approval_date,omitempty"`
	Admin         *Users         `gorm:"foreignKey:ApproverID;references:ID" json:"-"`
	ApproverID    *string        `gorm:"" json:"approver_id,omitempty"`
	RequestType   string         `gorm:"" json:"request_type,omitempty"`
}

type IssueRegistry struct {
	IssueID            string         `gorm:"primaryKey"`
	BookInventory      *BookInventory `gorm:"foreignKey:BookID;references:ISBN"`
	BookID             string         `gorm:""`
	Reader             *Users         `gorm:"foreignKey:ReaderID;references:ID"`
	ReaderID           string         `gorm:""`
	AdminIssue         *Users         `gorm:"foreignKey:IssueApproverID;references:ID"`
	IssueApproverID    string         `gorm:""`
	IssueStatus        string         `gorm:""`
	IssueDate          string         `gorm:""`
	ExpectedReturnDate string         `gorm:""`
	ReturnDate         *string        `gorm:""`
	AdminReturn        *Users         `gorm:"foreignKey:ReturnApproverID;references:ID"`
	ReturnApproverID   *string        `gorm:""`
}

type LibraryDetails struct {
	Library
	OwnerName  string `json:"owner_name"`
	OwnerEmail string `json:"owner_email"`
	TotalBooks int    `json:"total_books"`
}

type UpdateFields struct {
	Title     string `gorm:"column:title"`
	Authors   string `gorm:"column:authors"`
	Publisher string `gorm:"column:publisher"`
	Version   string `gorm:"column:version"`
}

type IssueRequestDetails struct {
	RequestEvents
	BookTitle       string `json:"book_title" binding:"required"`
	AvailableCopies int    `json:"available_copies" binding:"required"`
	ReaderName      string `json:"reader_name" binding:"required"`
}
