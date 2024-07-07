package handlers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"hng_stage_two_task/config"
	_ "hng_stage_two_task/docs"
	"hng_stage_two_task/internal/api/rest"
	"hng_stage_two_task/internal/dto"
	"hng_stage_two_task/internal/helper"
	"hng_stage_two_task/internal/repository"
	"hng_stage_two_task/internal/service"
	"hng_stage_two_task/internal/utils"

	"log"
	"net/http"
	"reflect"
)

type UserHandler struct {
	// svc UserService
	svc      service.UserService
	validate *validator.Validate
}

// SetupUserRoutes @Golang User/Organization API
// @version 1.0
// @description This is hng stage two task.
// @termsOfService http://swagger.io/terms/
// @contact.name @DevOlajide
// @contact.email programmerolajide@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8098
// @BasePath /
// SetupUserRoutes @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format: Bearer <token>
func SetupUserRoutes(rh *rest.RestHandler) {

	app := rh.App

	// create an instance of user service & inject to handler
	svc := service.UserService{
		Repo:    repository.NewUserRepository(rh.DB),
		OrgRepo: repository.NewOrganizationRepository(rh.DB),
		Auth:    rh.Auth,
		Config:  rh.Config,
	}

	handler := UserHandler{
		svc:      svc,
		validate: validator.New(),
	}

	app.Get("/", handler.Home)
	// Serve Swagger UI
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	publicRoutes := app.Group("/auth")

	// Public endpoints
	publicRoutes.Post("/register", handler.Register)
	publicRoutes.Post("/login", handler.Login)

	privateRoutes := app.Group("/api/", rh.Auth.Authorize)

	// Protected endpoints
	privateRoutes.Get("/users/:id", handler.GetUser)
	privateRoutes.Get("/organisations", handler.GetOrganisations)
	privateRoutes.Get("/organisations/:orgId", handler.GetOrganisationByOrgId)
	privateRoutes.Post("/organisations", handler.CreateOrganisation)
	privateRoutes.Post("/organisations/:orgId/users", handler.AddUserToOrganisation)

}

// Home godoc
// @Summary Home endpoint
// @Description Returns a hello world message
// @Tags Home
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router / [get]
func (h *UserHandler) Home(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"home": "Hello world",
	})
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with default organisation
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.UserSignupRequestDto true "User Signup Request"
// @Success 201 {object} dto.DefaultApiResponse
// @Failure 400 {object} dto.ErrorResponseDto
// @Failure 422 {object} dto.ErrorResponseDto
// @Router /auth/register [post]
func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	user := dto.UserSignupRequestDto{}
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status":     "Bad request",
			"message":    "Registration unsuccessful",
			"statusCode": http.StatusBadRequest,
		})
	}

	//Do data validations
	err := h.validate.Struct(user)
	if err != nil {
		// Handle validation errors
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errors []map[string]string

			// Use reflection to get custom error messages
			userType := reflect.TypeOf(user)
			for _, validationErr := range validationErrors {
				field := validationErr.StructField()
				message := validationErr.Tag() + " " + validationErr.Param()

				// Retrieve custom error message from struct tag
				if field, ok := userType.FieldByName(field); ok {
					customMessage := field.Tag.Get("error")
					if customMessage != "" {
						message = customMessage
					}
				}

				errors = append(errors, map[string]string{
					"field":   validationErr.Field(),
					"message": message,
				})
			}

			return ctx.Status(http.StatusUnprocessableEntity).JSON(map[string]interface{}{
				"errors": errors,
			})
		}

		// Handle other unexpected validation errors
		return ctx.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status":     "Bad request",
			"message":    "Registration unsuccessful",
			"statusCode": http.StatusBadRequest,
		})
	}

	// Check if user already exists
	exists, err := h.svc.UserExists(user.Email)
	if err != nil {
		log.Printf("Could not hceck if user with email: %v exist", user.Email)
		return ctx.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status":     "Bad request",
			"message":    "Registration unsuccessful",
			"statusCode": http.StatusBadRequest,
		})
	}

	if exists {
		return ctx.Status(http.StatusConflict).JSON(map[string]interface{}{
			"errors": []map[string]string{
				{
					"field":   "email",
					"message": fmt.Sprintf("User with email: %v already exists", user.Email),
				},
			},
		})
	}

	signupResponseData, err := h.svc.Signup(user)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status":     "Bad request",
			"message":    "Registration unsuccessful",
			"statusCode": http.StatusBadRequest,
		})
	}

	response := dto.DefaultApiResponse{
		BaseResponse: dto.BaseResponse[any]{
			Status:  "success",
			Message: "Registration successful",
			Data:    signupResponseData,
		},
	}
	return ctx.Status(http.StatusCreated).JSON(response)
}

// Login godoc
// @Summary Login a user
// @Description Login a user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.UserLoginDto true "User Login Request"
// @Success 200 {object} dto.DefaultApiResponse
// @Failure 400 {object} dto.ErrorResponseDto
// @Failure 401 {object} dto.ErrorResponseDto
// @Router /auth/login [post]
func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	loginDto := dto.UserLoginDto{}

	if err := ctx.BodyParser(&loginDto); err != nil {
		return helper.RespondWithError(ctx, http.StatusBadRequest, config.INVALID_PAYLOAD.Code, config.INVALID_PAYLOAD.Description)
	}

	//Do data validations
	err := h.validate.Struct(loginDto)
	if err != nil {
		// Handle validation errors
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errors []map[string]string

			// Use reflection to get custom error messages
			userType := reflect.TypeOf(loginDto)
			for _, validationErr := range validationErrors {
				field := validationErr.StructField()
				message := validationErr.Tag() + " " + validationErr.Param()

				// Retrieve custom error message from struct tag
				if field, ok := userType.FieldByName(field); ok {
					customMessage := field.Tag.Get("error")
					if customMessage != "" {
						message = customMessage
					}
				}

				errors = append(errors, map[string]string{
					"field":   validationErr.Field(),
					"message": message,
				})
			}

			return ctx.Status(http.StatusUnprocessableEntity).JSON(map[string]interface{}{
				"errors": errors,
			})
		}

		// Handle other unexpected validation errors
		return ctx.Status(http.StatusUnauthorized).JSON(map[string]interface{}{
			"status":     "Bad request",
			"message":    "Authentication failed",
			"statusCode": http.StatusUnauthorized,
		})
	}

	accessToken, err := h.svc.Login(loginDto.Email, loginDto.Password)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(map[string]interface{}{
			"status":     "Bad request",
			"message":    "Authentication failed",
			"statusCode": http.StatusUnauthorized,
		})
	}

	user, err := h.svc.Repo.FindUser(loginDto.Email)
	if err != nil {
		log.Printf("User with email: %v not found", loginDto.Email)
		return ctx.Status(http.StatusUnauthorized).JSON(map[string]interface{}{
			"status":     "Bad request",
			"message":    "Authentication failed",
			"statusCode": http.StatusUnauthorized,
		})
	}

	loginResponseData := dto.LoginResponse{
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

	response := dto.DefaultApiResponse{
		BaseResponse: dto.BaseResponse[any]{
			Status:  "success",
			Message: "Login successful",
			Data:    loginResponseData,
		},
	}
	return ctx.Status(http.StatusOK).JSON(response)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get user information by user ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.DefaultApiResponse
// @Failure 400 {object} dto.ErrorResponseDto
// @Failure 404 {object} dto.ErrorResponseDto
// @Security BearerAuth
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(ctx *fiber.Ctx) error {
	// Extract the id from params
	userID := ctx.Params("id")
	if utils.IsEmpty(userID) {
		return ctx.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status":     "Bad request",
			"message":    "User id is required",
			"statusCode": http.StatusBadRequest,
		})
	}

	// Get the current logged-in user from token
	user := h.svc.Auth.GetCurrentUser(ctx)

	// Ensure the user is properly set
	if user.UserID == "" {
		return ctx.Status(http.StatusUnauthorized).JSON(map[string]interface{}{
			"status":     "Bad request",
			"message":    "Unauthorized to access resource",
			"statusCode": http.StatusUnauthorized,
		})
	}

	// Call service to get the user by ID
	foundUser, err := h.svc.GetUserById(userID, user)
	if err != nil {
		var statusCode int
		var message string
		switch err.Error() {
		case "unauthorized access":
			statusCode = http.StatusUnauthorized
			message = "Unauthorized to access resource"
		case "user not found":
			statusCode = http.StatusNotFound
			message = "User not found"
		default:
			statusCode = http.StatusBadRequest
			message = "Client error"
		}
		errorResponse := dto.ErrorResponseDto{
			Status:     "Bad request",
			Message:    message,
			StatusCode: statusCode,
		}
		return ctx.Status(statusCode).JSON(errorResponse)
	}

	// Prepare the response data
	responseData := dto.FindUserByIdResponseData{
		UserID:    foundUser.UserID,
		FirstName: foundUser.FirstName,
		LastName:  foundUser.LastName,
		Email:     foundUser.Email,
		Phone:     foundUser.Phone,
	}

	response := dto.DefaultApiResponse{
		BaseResponse: dto.BaseResponse[any]{
			Status:  "success",
			Message: "User retrieved successfully",
			Data:    responseData,
		},
	}
	return ctx.Status(http.StatusOK).JSON(response)
}

// GetOrganisations godoc
// @Summary Get all organisations
// @Description Get a list of all organisations
// @Tags Organisation
// @Accept json
// @Produce json
// @Success 200 {array} dto.OrganisationsResponse
// @Failure 400 {object} dto.ErrorResponseDto
// @Security BearerAuth
// @Router /api/organisations [get]
func (h *UserHandler) GetOrganisations(ctx *fiber.Ctx) error {
	// Get the current logged-in user from token
	user := h.svc.Auth.GetCurrentUser(ctx)

	// Call service to get the user's organisations
	organisations, err := h.svc.GetUserOrganisations(user.UserID)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			dto.ErrorResponseDto{
				Status:     "Bad request",
				Message:    "Client error",
				StatusCode: http.StatusBadRequest,
			})
	}

	responseData := dto.OrganisationsResponse{
		Organisations: organisations,
	}

	response := dto.DefaultApiResponse{
		BaseResponse: dto.BaseResponse[any]{
			Status:  "success",
			Message: "Users organizations retrieved successfully",
			Data:    responseData,
		},
	}
	return ctx.Status(http.StatusOK).JSON(response)
}

// GetOrganisationByOrgId godoc
// @Summary Get organisation by ID
// @Description Get organisation information by organisation ID
// @Tags Organisation
// @Accept json
// @Produce json
// @Param orgId path string true "Organisation ID"
// @Success 200 {object} dto.OrganisationsResponse
// @Failure 400 {object} dto.ErrorResponseDto
// @Failure 404 {object} dto.ErrorResponseDto
// @Security BearerAuth
// @Router /api/organisations/{orgId} [get]
func (h *UserHandler) GetOrganisationByOrgId(ctx *fiber.Ctx) error {
	// Extract the orgId from params
	orgID := ctx.Params("orgId")
	log.Printf("Extracted orgID: %s", orgID) // Debug logging

	// Get the current logged-in user from token
	user := h.svc.Auth.GetCurrentUser(ctx)
	log.Printf("The current user is: %v", user)

	// Call service to get the organisation by ID
	organisation, err := h.svc.GetOrganisationById(user.UserID, orgID)
	if err != nil {
		if err.Error() == "organisation not found" {
			return ctx.Status(http.StatusNotFound).JSON(map[string]interface{}{
				"status":     "Bad request",
				"message":    "Organisation not found",
				"statusCode": http.StatusNotFound,
			})
		}

		if err.Error() == "user not found" {
			return ctx.Status(http.StatusUnauthorized).JSON(map[string]interface{}{
				"status":     "Bad request",
				"message":    "Unauthorized to access resource",
				"statusCode": http.StatusUnauthorized,
			})
		}

		return ctx.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status":     "Bad request",
			"message":    "Client error",
			"statusCode": http.StatusBadRequest,
		})
	}

	response := dto.DefaultApiResponse{
		BaseResponse: dto.BaseResponse[any]{
			Status:  "success",
			Message: "Organization retrieved successfully",
			Data:    organisation,
		},
	}
	return ctx.Status(http.StatusOK).JSON(response)
}

// CreateOrganisation godoc
// @Summary Create a new organisation
// @Description Create a new organisation
// @Tags Organisation
// @Accept json
// @Produce json
// @Param organisation body dto.CreateOrganisationRequest true "Create Organisation Request"
// @Success 201 {object} dto.DefaultApiResponse
// @Failure 400 {object} dto.ErrorResponseDto
// @Failure 422 {object} dto.ErrorResponseDto
// @Security BearerAuth
// @Router /api/organisations [post]
func (h *UserHandler) CreateOrganisation(ctx *fiber.Ctx) error {
	// Get the current logged-in user from token
	user := h.svc.Auth.GetCurrentUser(ctx)

	// Ensure the user is properly set
	if user.UserID == "" {
		return ctx.Status(http.StatusUnauthorized).JSON(map[string]interface{}{
			"status":     "Bad request",
			"message":    "Unauthorized to access resource",
			"statusCode": http.StatusUnauthorized,
		})
	}

	// Parse and validate the request body
	var req dto.CreateOrganisationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status":     "Bad Request",
			"message":    "Client error",
			"statusCode": http.StatusBadRequest,
		})
	}

	// Do data validations for the request body
	err := h.validate.Struct(req)
	if err != nil {
		// Handle validation errors
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errors []map[string]string

			// Use reflection to get custom error messages
			reqType := reflect.TypeOf(req)
			for _, validationErr := range validationErrors {
				field := validationErr.StructField()
				message := validationErr.Tag() + " " + validationErr.Param()

				// Retrieve custom error message from struct tag
				if field, ok := reqType.FieldByName(field); ok {
					customMessage := field.Tag.Get("error")
					if customMessage != "" {
						message = customMessage
					}
				}

				errors = append(errors, map[string]string{
					"field":   validationErr.Field(),
					"message": message,
				})
			}

			return ctx.Status(http.StatusUnprocessableEntity).JSON(map[string]interface{}{
				"errors": errors,
			})
		}

		// Handle other unexpected validation errors
		return ctx.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status":     "Bad request",
			"message":    "Validation unsuccessful",
			"statusCode": http.StatusBadRequest,
		})
	}

	// Call service to create the organisation
	organisation, err := h.svc.CreateOrganisation(user.UserID, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(map[string]interface{}{
			"status":     "Internal server error",
			"message":    "Failed to create organisation",
			"statusCode": http.StatusInternalServerError,
		})
	}

	response := dto.DefaultApiResponse{
		BaseResponse: dto.BaseResponse[any]{
			Status:  "success",
			Message: "Organisation created successfully",
			Data:    organisation,
		},
	}
	return ctx.Status(http.StatusCreated).JSON(response)
}

// AddUserToOrganisation godoc
// @Summary Add a user to an organisation
// @Description Add a user to an organisation by organisation ID and user ID
// @Tags Organisation
// @Accept json
// @Produce json
// @Param orgId path string true "Organisation ID"
// @Param user body dto.AddUserToOrganisationRequestDto true "Add User to Organisation Request"
// @Success 201 {object} dto.DefaultApiResponse
// @Failure 400 {object} dto.ErrorResponseDto
// @Failure 422 {object} dto.ErrorResponseDto
// @Security BearerAuth
// @Router /api/organisations/{orgId}/users [post]
func (h *UserHandler) AddUserToOrganisation(ctx *fiber.Ctx) error {
	// Get the organisation ID from the URL parameters
	orgId := ctx.Params("orgId")

	// Get the current logged-in user from token
	user := h.svc.Auth.GetCurrentUser(ctx)
	// Ensure the user is properly set
	if user.UserID == "" {
		return ctx.Status(http.StatusUnauthorized).JSON(map[string]interface{}{
			"status":     "Bad request",
			"message":    "Unauthorized to access resource",
			"statusCode": http.StatusUnauthorized,
		})
	}

	if orgId == "" {
		return ctx.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status":     "Bad request",
			"message":    "Organisation ID is required",
			"statusCode": http.StatusBadRequest,
		})
	}

	// Parse and validate the request body
	var req dto.AddUserToOrganisationRequestDto
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status":     "Bad Request",
			"message":    "Client error",
			"statusCode": http.StatusBadRequest,
		})
	}

	// Validate the request data
	err := h.validate.Struct(req)
	if err != nil {
		// Handle validation errors
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errors []map[string]string

			// Use reflection to get custom error messages
			reqType := reflect.TypeOf(req)
			for _, validationErr := range validationErrors {
				field := validationErr.StructField()
				message := validationErr.Tag() + " " + validationErr.Param()

				// Retrieve custom error message from struct tag
				if field, ok := reqType.FieldByName(field); ok {
					customMessage := field.Tag.Get("error")
					if customMessage != "" {
						message = customMessage
					}
				}

				errors = append(errors, map[string]string{
					"field":   validationErr.Field(),
					"message": message,
				})
			}

			return ctx.Status(http.StatusUnprocessableEntity).JSON(map[string]interface{}{
				"errors": errors,
			})
		}

		// Handle other unexpected validation errors
		return ctx.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status":     "Bad request",
			"message":    "Validation unsuccessful",
			"statusCode": http.StatusBadRequest,
		})
	}

	// Call service to add the user to the organisation
	err = h.svc.AddUserToOrganisation(orgId, req.UserID, user)
	if err != nil {
		var statusCode int
		var message string
		switch err.Error() {
		case "unauthorized access":
			statusCode = http.StatusUnauthorized
			message = "Unauthorized to access resource"
		case "user not found":
			statusCode = http.StatusNotFound
			message = "User not found"
		case "LoggedIn user is not authorized to add users to this organization":
			statusCode = http.StatusUnauthorized
			message = "LoggedIn user is not authorized to add users to this organization"
		default:
			statusCode = http.StatusBadRequest
			message = "Client error"
		}
		errorResponse := dto.ErrorResponseDto{
			Status:     "Bad request",
			Message:    message,
			StatusCode: statusCode,
		}
		return ctx.Status(statusCode).JSON(errorResponse)
	}

	// Return success response
	return ctx.Status(http.StatusOK).JSON(map[string]interface{}{
		"status":  "success",
		"message": "User added to organisation successfully",
	})
}
