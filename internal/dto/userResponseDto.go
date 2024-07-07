package dto

type UserResponse struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type AuthResponseData struct {
	AccessToken string       `json:"accessToken"`
	User        UserResponse `json:"user"`
}

type LoginResponse struct {
	AuthResponseData
}

type SignupResponse struct {
	AuthResponseData
}

type FindUserByIdResponseData struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
