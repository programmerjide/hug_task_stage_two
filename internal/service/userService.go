package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"hng_stage_two_task/config"
	"hng_stage_two_task/internal/domain"
	"hng_stage_two_task/internal/dto"
	"hng_stage_two_task/internal/helper"
	"hng_stage_two_task/internal/repository"
	"log"
)

type UserService struct {
	Repo    repository.UserRepository
	OrgRepo repository.OrganizationRepository
	Auth    helper.Auth
	Config  config.AppConfig
}

func (us UserService) Signup(createUserRequestDTO dto.UserSignupRequestDto) (dto.SignupResponse, error) {
	log.Printf("User signup with email: %v", createUserRequestDTO.Email)

	hashPassword, err := us.Auth.CreateHashedPassword(createUserRequestDTO.Password)
	if err != nil {
		return dto.SignupResponse{}, err
	}

	user, err := us.Repo.CreateUser(domain.User{
		FirstName: createUserRequestDTO.FirstName,
		LastName:  createUserRequestDTO.LastName,
		Email:     createUserRequestDTO.Email,
		Password:  hashPassword,
		Phone:     createUserRequestDTO.Phone,
	})
	if err != nil {
		return dto.SignupResponse{}, err
	}

	// Create the organisation
	organisation := domain.Organisation{
		OrgID:       uuid.New().String(),
		Name:        fmt.Sprintf("%s's Organisation", user.FirstName),
		Description: fmt.Sprintf("Default organization for %s %s", user.FirstName, user.LastName),
	}

	// Create default organization
	createdOrganisation, err := us.Repo.CreateOrganisation(organisation)
	if err != nil {
		return dto.SignupResponse{}, err
	}

	// Link user to organisation
	err = us.Repo.AddUserToOrganisation(user.UserID, createdOrganisation.OrgID)
	if err != nil {
		log.Printf("Error linking user to organisation: %v", err)
		return dto.SignupResponse{}, err
	}

	// Update user object with organizations
	user.Orgs = append(user.Orgs, createdOrganisation)

	// Generate access token with updated user object
	accessToken, err := us.Auth.GenerateAccessToken(user)
	if err != nil {
		log.Printf("error generating access token with errors: %v", err)
		return dto.SignupResponse{}, err
	}

	// Create response
	signResponseData := dto.SignupResponse{
		AuthResponseData: dto.AuthResponseData{
			AccessToken: accessToken,
			User: dto.UserResponse{
				UserID:    user.UserID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Phone:     user.Phone,
			},
		},
	}

	return signResponseData, nil
}

func (us UserService) findUserByEmail(email string) (*domain.User, error) {
	// perform db  operation
	user, err := us.Repo.FindUser(email)
	return &user, err
}

func (us UserService) Login(email string, password string) (string, error) {

	user, err := us.findUserByEmail(email)
	if err != nil {
		return "", errors.New("user does not exist with email")
	}

	err = us.Auth.VerifyPassword(password, user.Password)
	if err != nil {
		return "", err
	}

	// generate token
	return us.Auth.GenerateAccessToken(*user)
}

func (us UserService) UserExists(email string) (bool, error) {
	user, err := us.findUserByEmail(email)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

func (us UserService) GetUserById(userID string, requestingUser domain.User) (dto.FindUserByIdResponseData, error) {
	// Fetch the user data by user ID to check if the user exists
	user, err := us.Repo.FindUserById(userID)
	if err != nil {
		return dto.FindUserByIdResponseData{}, fmt.Errorf("user not found: %w", err)
	}

	// Fetch all organization IDs for both requesting user and requested user
	requestingUserOrgIDs, err := us.GetOrganizationIDsByUserID(requestingUser.UserID)
	if err != nil {
		return dto.FindUserByIdResponseData{}, fmt.Errorf("failed to get organization IDs for requesting user: %w", err)
	}

	requestedUserOrgIDs, err := us.GetOrganizationIDsByUserID(userID)
	if err != nil {
		return dto.FindUserByIdResponseData{}, fmt.Errorf("failed to get organization IDs for requested user: %w", err)
	}

	// Check if there's a common organization between the two users
	if !us.isUserInSameOrganization(requestingUserOrgIDs, requestedUserOrgIDs) {
		return dto.FindUserByIdResponseData{}, fmt.Errorf("user not in organization")
	}

	// Prepare the response data
	responseData := dto.FindUserByIdResponseData{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
	}

	return responseData, nil
}

func (us UserService) GetOrganizationIDsByUserID(userID string) ([]string, error) {
	user, err := us.Repo.FindUserById(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	orgIDs := make([]string, len(user.Orgs))
	for i, org := range user.Orgs {
		orgIDs[i] = org.OrgID
	}

	return orgIDs, nil
}

func (us UserService) isUserInSameOrganization(requestingUserOrgIDs, requestedUserOrgIDs []string) bool {
	orgIDMap := make(map[string]struct{})
	for _, orgID := range requestingUserOrgIDs {
		orgIDMap[orgID] = struct{}{}
	}

	for _, orgID := range requestedUserOrgIDs {
		if _, exists := orgIDMap[orgID]; exists {
			return true
		}
	}

	return false
}

func (us UserService) GetUserOrganisations(userID string) ([]dto.OrganisationResponse, error) {
	user, err := us.Repo.FindUserById(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	var organisations []dto.OrganisationResponse
	for _, org := range user.Orgs {
		organisations = append(organisations, dto.OrganisationResponse{
			OrgID:       org.OrgID,
			Name:        org.Name,
			Description: org.Description,
		})
	}

	return organisations, nil
}

func (us UserService) GetOrganisationById(userID, orgID string) (dto.OrganisationResponse, error) {
	user, err := us.Repo.FindUserById(userID)
	if err != nil {
		return dto.OrganisationResponse{}, fmt.Errorf("user not found: %w", err)
	}

	var organisation domain.Organisation
	found := false
	for _, org := range user.Orgs {
		if org.OrgID == orgID {
			organisation = org
			found = true
			break
		}
	}

	if !found {
		return dto.OrganisationResponse{}, fmt.Errorf("organisation not found")
	}

	return dto.OrganisationResponse{
		OrgID:       organisation.OrgID,
		Name:        organisation.Name,
		Description: organisation.Description,
	}, nil
}

func (us UserService) CreateOrganisation(userID string, req dto.CreateOrganisationRequest) (dto.OrganisationResponse, error) {
	// Create a new organisation instance
	organisation := domain.Organisation{
		Name:        req.Name,
		Description: req.Description,
	}

	// Save the organisation to the database
	createdOrg, err := us.Repo.CreateOrganisation(organisation)
	if err != nil {
		return dto.OrganisationResponse{}, fmt.Errorf("failed to create organisation: %w", err)
	}

	// Add the organisation to the user's organisations
	err = us.Repo.AddUserToOrganisation(userID, createdOrg.OrgID)
	if err != nil {
		return dto.OrganisationResponse{}, fmt.Errorf("failed to add user to organisation: %w", err)
	}

	return dto.OrganisationResponse{
		OrgID:       createdOrg.OrgID,
		Name:        createdOrg.Name,
		Description: createdOrg.Description,
	}, nil
}

func (us UserService) AddUserToOrganisation(orgId, userId string, currentUser domain.User) error {
	if us.Repo == nil || us.OrgRepo == nil {
		return fmt.Errorf("repository is not initialized")
	}

	// Logic to add the user to the organisation
	_, err := us.Repo.FindUserById(userId)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	_, err = us.OrgRepo.FindOrganizationById(orgId)
	if err != nil {
		return fmt.Errorf("organization does not exist: %w", err)
	}

	isInOrg, err := us.Repo.IsUserInOrganisation(currentUser.UserID, orgId)
	if err != nil {
		return fmt.Errorf("failed to check user organization association: %w", err)
	}
	if !isInOrg {
		return fmt.Errorf("LoggedIn user is not authorized to add users to this organization")
	}

	// and then updating the database to link the user to the organisation.
	err = us.Repo.AddUserToOrganisation(userId, orgId)
	if err != nil {
		return fmt.Errorf("failed to add user to organization: %w", err)
	}

	return nil
}
