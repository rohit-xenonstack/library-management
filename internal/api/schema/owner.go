package schema

import "library-management/backend/internal/api/model"

type CreateLibraryRequest struct {
	LibraryName   string `json:"library_name" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required"`
	ContactNumber string `json:"contact" binding:"required"`
}

type CreateAdminRequest struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required"`
	ContactNumber string `json:"contact" binding:"required"`
	LibID         string `json:"library_id" binding:"required"`
}

type GetLibrariesResponse struct {
	RequiredResponseFields
	Libraries *[]model.LibraryDetails `json:"libraries,omitempty"`
}

type GetAdminsRequest struct {
	LibID string `json:"library_id" binding:"required"`
}
type GetAdminsResponse struct {
	RequiredResponseFields
	Admins *[]model.Users `json:"admins,omitempty"`
}
