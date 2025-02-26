package schema

type RaiseIssueRequest struct {
	ReaderEmail string `json:"email" binding:"required"`
	BookID      string `json:"isbn" binding:"required"`
}

type GetLatestAvailabilityResponse struct {
	RequiredResponseFields
	Date *string `json:"date,omitempty"`
}
