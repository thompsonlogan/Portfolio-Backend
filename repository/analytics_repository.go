package repository

import (
	"api/internal/data/model"

	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	GetVisitByIpAndSource(ip string, src string) (*model.PortfolioVisit, error)
	AddVisit(visit *model.PortfolioVisit) error
  UpdateVisit(visit *model.PortfolioVisit) error
}

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) GetVisitByIpAndSource(ip string, source string) (*model.PortfolioVisit, error) {
	var visit model.PortfolioVisit
	if err := r.db.Where("ip = ? AND source = ?", ip, source).First(&visit).Error; err != nil {
		return nil, err
	}
	return &visit, nil
}

func (r *analyticsRepository) AddVisit(visit *model.PortfolioVisit) error {
  return r.db.
    Where("ip = ? AND source = ?", visit.IP, visit.Source).
    FirstOrCreate(visit).
    Error
}

func (r *analyticsRepository) UpdateVisit(visit *model.PortfolioVisit) error {
	return r.db.Model(visit).Updates(visit).Error
}