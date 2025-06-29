package service

import (
	"fmt"
	"payment/internal/model"
	"payment/internal/repository"
)

type PackageService struct {
	packageRepository *repository.PackageRepository
}

func NewPackageService(repo *repository.PackageRepository) *PackageService {
	return &PackageService{packageRepository: repo}
}

func (s *PackageService) GetAllPackages() ([]model.Packages, error) {
	packages, err := s.packageRepository.FindAllPackage()
	if err != nil {
		return nil, fmt.Errorf("failed to get packages: %w", err)
	}
	return packages, nil

}

func (s *PackageService) GetPackageById(id int) (*model.Packages, error) {
	packages, err := s.packageRepository.FindPackageById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get packages: %w", err)
	}
	return packages, nil

}
