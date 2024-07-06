package repository

import (
	"gorm.io/gorm"
	"hng_stage_two_task/internal/domain"
)

type OrganizationRepository interface {
	FindOrganizationById(orgId string) (domain.Organisation, error)
}

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{
		db: db,
	}
}

func (r organizationRepository) FindOrganizationById(orgId string) (domain.Organisation, error) {
	var organization domain.Organisation
	if err := r.db.Where("org_id = ?", orgId).First(&organization).Error; err != nil {
		return domain.Organisation{}, err
	}
	return organization, nil
}
