package schema

import "library-management/backend/internal/api/model"

type AddBookRequest struct {
	AdminEmail string `json:"email" binding:"required"`
	ISBN       string `json:"isbn" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Authors    string `json:"authors" binding:"required"`
	Publisher  string `json:"publisher" binding:"required"`
	Version    string `json:"version" binding:"required"`
}

type RemoveBookRequest struct {
	ISBN string `json:"isbn" binding:"required"`
}

type ListIssueRequestResponse struct {
	RequiredResponseFields
	Requests []model.IssueRequestDetails `json:"requests"`
}

type RequestDetails struct {
	RequestID string `json:"request_id" binding:"required"`
	AdminID   string `json:"user_id" binding:"required"`
}

type UpdateBookRequest struct {
	ISBN      string `json:"isbn" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Authors   string `json:"authors" binding:"required"`
	Publisher string `json:"publisher" binding:"required"`
	Version   string `json:"version" binding:"required"`
}

type GetBookByISBNResponse struct {
	Status  string               `json:"status"`
	Payload *model.BookInventory `json:"payload"`
}
