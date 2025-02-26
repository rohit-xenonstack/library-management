package schema

import "library-management/backend/internal/api/model"

type RequiredResponseFields struct {
	Status  string `json:"status" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required"`
}
type LoginResponse struct {
	RequiredResponseFields
	AccessToken *string      `json:"access_token,omitempty"`
	User        *model.Users `json:"user,omitempty"`
}

type UserDetailsResponse struct {
	RequiredResponseFields
	User *model.Users `json:"user" binding:"required"`
}

type ReaderSignupRequest struct {
	LibID   string `json:"library_id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required"`
	Contact string `json:"contact" binding:"required"`
}
type ReaderSignupResponse struct {
	RequiredResponseFields
}

type CreateLibraryWithOwnerRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Contact     string `json:"contact" binding:"required"`
	LibraryName string `json:"library_name" binding:"required"`
}

type RefreshAccessTokenResponse struct {
	RequiredResponseFields
	AccessToken *string `json:"access_token" binding:"required"`
}
