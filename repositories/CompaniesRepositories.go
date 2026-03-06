package repositories

import (
	"Gin/initializers"
	"Gin/models"
	"context"
	"time"
)

type ICompaniesRepository interface {
	PostCompany(user *models.Company) error
	GetCompanies() ([]models.Company, error)
	GetCompany(companyID uint) (*models.Company, error)
}

type companiesRepository struct{}

func NewCompaniesRepository() ICompaniesRepository {
	return &companiesRepository{}
}

func (r *companiesRepository) PostCompany(user *models.Company) error {
	result := initializers.DB.Create(user)
	return result.Error
}

func (r *companiesRepository) GetCompanies() ([]models.Company, error) {
	context, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var companies []models.Company
	result := initializers.DB.WithContext(context).Find(&companies)
	if result.Error != nil {
		return nil, result.Error
	}
	return companies, nil
}

func (r *companiesRepository) GetCompany(companyID uint) (*models.Company, error) {
	context, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var company models.Company
	result := initializers.DB.WithContext(context).Where("id = ?", companyID).First(&company)
	if result.Error != nil {
		return nil, result.Error
	}
	return &company, nil
}
