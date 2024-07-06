package dto

type CreateOrganisationRequest struct {
	Name        string `json:"name" validate:"required" error:"Organization name is required"`
	Description string `json:"description"`
}
