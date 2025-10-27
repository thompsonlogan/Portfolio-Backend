package services

import (
	"api/internal/data/model"
	"api/repository"
	"errors"

	"gorm.io/gorm"
)

type AnalyticsService interface {
	AddVisit(visit *model.PortfolioVisit) error
  AddGithubVisit(visit *model.PortfolioVisit) error
  AddLinkedinVisit(visit *model.PortfolioVisit) error
  AddResumeDownload(visit *model.PortfolioVisit) error
}

type analyticsService struct {
	repo repository.AnalyticsRepository
}

// Ternary is a helper function to mimic the ternary operator
func Ternary(condition bool, valueIfTrue, valueIfFalse interface{}) interface{} {
    if condition {
        return valueIfTrue
    }
    return valueIfFalse
}

func NewAnalyticsService(repo repository.AnalyticsRepository) AnalyticsService {
  return &analyticsService{repo: repo}
}

func (s *analyticsService)AddVisit(visit *model.PortfolioVisit) error {
  existingVisit, err := s.repo.GetVisitByIpAndSource(visit.IP, visit.Source)

  if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      visit.VisitCount = 1
      return s.repo.AddVisit(visit)
    }
    return err
  }

  existingVisit.VisitCount++
  return s.repo.UpdateVisit(existingVisit)
}

func (s *analyticsService)AddGithubVisit(visit *model.PortfolioVisit) error {
  existingVisit, err := s.repo.GetVisitByIpAndSource(visit.IP, visit.Source)

    if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      visit.VisitCount = 1
      visit.GithubVisitCount = 1
      return s.repo.AddVisit(visit)
    }
    return err
  }

  existingVisit.GithubVisitCount++
  return s.repo.UpdateVisit(existingVisit)
}

func (s *analyticsService)AddLinkedinVisit(visit *model.PortfolioVisit) error {
  existingVisit, err := s.repo.GetVisitByIpAndSource(visit.IP, visit.Source)

    if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      visit.VisitCount = 1
      visit.LinkedinVisitCount = 1
      return s.repo.AddVisit(visit)
    }
    return err
  }

  existingVisit.LinkedinVisitCount++
  return s.repo.UpdateVisit(existingVisit)
}

func (s *analyticsService)AddResumeDownload(visit *model.PortfolioVisit) error {
  existingVisit, err := s.repo.GetVisitByIpAndSource(visit.IP, visit.Source)

    if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      visit.VisitCount = 1
      visit.ResumeDownloadCount = 1
      return s.repo.AddVisit(visit)
    }
    return err
  }

  existingVisit.ResumeDownloadCount++
  return s.repo.UpdateVisit(existingVisit)
}