package schema

import "library-management/backend/internal/api/model"

type SearchBookRequest struct {
	SearchString string `json:"search_string" binding:"required"`
	Field        string `json:"field" binding:"required"`
}

type SearchBookResponse struct {
	RequiredResponseFields
	Books *[]model.BookInventory `json:"books,omitempty"`
}

type SearchBookByISBNResponse struct {
	RequiredResponseFields
	Book *model.BookInventory `json:"book,omitempty"`
}
