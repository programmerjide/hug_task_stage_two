package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hng_stage_two_task/internal/domain"
	"log"
)

type UserRepository interface {
	CreateUser(u domain.User) (domain.User, error)
	FindUser(email string) (domain.User, error)
	FindUserById(userID string) (domain.User, error)
	UpdateUser(id uint, u domain.User) (domain.User, error)
	CreateOrganisation(organisation domain.Organisation) (domain.Organisation, error)
	AddUserToOrganisation(userID, orgID string) error
	IsUserInOrganisation(userID, orgID string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r userRepository) CreateUser(user domain.User) (domain.User, error) {
	user.UserID = uuid.New().String()

	err := r.db.Create(&user).Error
	if err != nil {
		log.Printf(" Error occurred while creating user %v", err)
		return domain.User{}, errors.New("failed to create user")
	}

	return user, nil
}

func (r userRepository) FindUser(email string) (domain.User, error) {

	var user domain.User

	err := r.db.First(&user, "email=?", email).Error

	if err != nil {
		log.Printf(" Error occurred while finding user %v", err)
		return domain.User{}, errors.New("user does not exist")
	}

	return user, nil
}

func (r userRepository) FindUserById(userID string) (domain.User, error) {
	var user domain.User
	if err := r.db.Preload("Orgs").Where("user_id = ?", userID).First(&user).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r userRepository) UpdateUser(id uint, u domain.User) (domain.User, error) {

	var user domain.User

	err := r.db.Model(&user).Clauses(clause.Returning{}).Where("id=?", id).Updates(u).Error

	if err != nil {
		log.Printf("error on update %v", err)
		return domain.User{}, errors.New("failed update user")
	}

	return user, nil
}

func (r userRepository) CreateOrganisation(organisation domain.Organisation) (domain.Organisation, error) {
	organisation.OrgID = uuid.New().String()
	if err := r.db.Create(&organisation).Error; err != nil {
		return domain.Organisation{}, err
	}
	return organisation, nil
}

func (r userRepository) AddUserToOrganisation(userID, orgID string) error {
	var user domain.User
	err := r.db.First(&user, "user_id=?", userID).Error
	if err != nil {
		log.Printf("Error finding user by ID: %v", err)
		return err
	}

	var organisation domain.Organisation
	err = r.db.First(&organisation, "org_id=?", orgID).Error
	if err != nil {
		log.Printf("Error finding organisation by ID: %v", err)
		return err
	}

	user.Orgs = append(user.Orgs, organisation)

	err = r.db.Save(&user).Error
	if err != nil {
		log.Printf("Error saving user with new organisation: %v", err)
		return err
	}

	return nil
}

func (r userRepository) IsUserInOrganisation(userID, orgID string) (bool, error) {
	var user domain.User
	err := r.db.Preload("Orgs").First(&user, "user_id = ?", userID).Error
	if err != nil {
		return false, err
	}

	for _, org := range user.Orgs {
		if org.OrgID == orgID {
			return true, nil
		}
	}
	return false, nil
}
