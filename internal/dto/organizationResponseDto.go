package dto

type OrganisationResponse struct {
	OrgID       string `json:"orgId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OrganisationsResponse struct {
	Organisations []OrganisationResponse `json:"organisations"`
}
