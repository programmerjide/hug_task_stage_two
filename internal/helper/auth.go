package helper

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"hng_stage_two_task/internal/domain"
	"hng_stage_two_task/internal/utils"
	"log"
	"net/http"
	"strings"
	"time"
)

type Auth struct {
	Secret string
}

func SetupAuth(s string) Auth {
	return Auth{
		Secret: s,
	}
}

func (a Auth) CreateHashedPassword(password string) (string, error) {

	if len(password) < 6 {
		return "", errors.New("password length should be at least 6 characters long")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", errors.New("error generating hashedPassword")

	}

	return string(hashPassword), nil
}

func (a Auth) GenerateAccessToken(user domain.User) (string, error) {
	if utils.IsZero(user.UserID) || utils.IsEmpty(user.Email) {
		return "", errors.New("user_id and email are required to generate token")
	}

	// Collect organization IDs from user object
	orgIDs := make([]string, len(user.Orgs))
	for i, org := range user.Orgs {
		orgIDs[i] = org.OrgID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":           user.UserID,
		"firstName":        user.FirstName,
		"lastName":         user.LastName,
		"email":            user.Email,
		"phone":            user.Phone,
		"organisation_ids": orgIDs,
		"exp":              time.Now().Add(time.Hour * 1).Unix(),
	})

	accessToken, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		return "", errors.New("error generating token")
	}

	return accessToken, nil
}

func (a Auth) VerifyPassword(password string, hashPassword string) error {

	if len(password) < 5 {
		return errors.New("password length should be at least 5 characters long")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return errors.New("password does not match")
	}

	return nil

}

func (a Auth) VerifyAccessToken(token string) (domain.User, error) {
	tokenArr := strings.Split(token, " ")
	if len(tokenArr) != 2 {
		return domain.User{}, errors.New("invalid token format")
	}

	if tokenArr[0] != "Bearer" {
		return domain.User{}, errors.New("invalid token prefix")
	}

	tokenStr := tokenArr[1]

	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unknown signing method %v", t.Header)
		}
		return []byte(a.Secret), nil
	})

	if err != nil {
		return domain.User{}, errors.New("invalid signing method")
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return domain.User{}, errors.New("token is expired")
		}

		// Extract organization IDs from the claims and populate the Orgs field
		orgIDs := claims["organisation_ids"].([]interface{})
		var orgs []domain.Organisation
		for _, orgID := range orgIDs {
			orgs = append(orgs, domain.Organisation{
				OrgID: orgID.(string),
			})
		}

		user := domain.User{
			UserID:    claims["userId"].(string),
			Email:     claims["email"].(string),
			FirstName: claims["firstName"].(string),
			LastName:  claims["lastName"].(string),
			Phone:     claims["phone"].(string),
			Orgs:      orgs,
		}

		return user, nil
	}

	return domain.User{}, errors.New("token verification failed")
}

func (a Auth) Authorize(ctx *fiber.Ctx) error {
	// Retrieve Authorization header
	authHeader := ctx.Get("Authorization")

	// Verify access token
	user, err := a.VerifyAccessToken(authHeader)
	if err != nil {
		// Handle any error
		log.Printf("Token verification error: %v", err)
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Authorization failed",
		})
	}

	// Access granted, store user in context
	ctx.Locals("user", user)
	return ctx.Next()
}

func (a Auth) GetCurrentUser(ctx *fiber.Ctx) domain.User {
	user := ctx.Locals("user")
	return user.(domain.User)
}
