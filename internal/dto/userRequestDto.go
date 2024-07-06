package dto

type UserLoginDto struct {
	Email    string `json:"email" validate:"required,email" error:"Email is required and must be a valid email address"`
	Password string `json:"password" validate:"required" error:"Password is required"`
}

type UserSignupRequestDto struct {
	UserLoginDto
	Phone     string `json:"phone"`
	FirstName string `json:"firstName" validate:"required,min=3,max=50" error:"First name is required and must be between 3 to 50 characters"`
	LastName  string `json:"lastName" validate:"required,min=3,max=50" error:"First name is required and must be between 3 to 50 characters"`
}

type AddUserToOrganisationRequestDto struct {
	UserID string `json:"userId" validate:"required" error:"UserID is required"`
}
