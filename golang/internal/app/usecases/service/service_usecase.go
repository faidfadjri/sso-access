package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"

	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/repositories"
	"akastra-access/internal/infrastructure/utils/image"
	"akastra-access/internal/infrastructure/utils/pagination"
	"akastra-access/internal/interface/http/request"
	pkgErrors "akastra-access/internal/pkg/errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ServiceUsecase handles HTTP requests
type serviceUsecase struct {
	serviceRepository repositories.ServiceRepository
}

func NewServiceUsecase(
	repo repositories.ServiceRepository,
) ServiceUsecase {
	return &serviceUsecase{
		serviceRepository: repo,
	}
}

func (t *serviceUsecase) CreateClient(req request.CreateClientReq, file multipart.File, header *multipart.FileHeader) (*domain.Services, error) {
	// Process image
	logoPath, err := image.ProcessImage(file, header, "public/images/services")
	if err != nil {
		return nil, err
	}

	clientID, err := generateClientID()
	if err != nil {
		return nil, err
	}

	clientSecret, err := generateClientSecret()
	if err != nil {
		return nil, err
	}

	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create service object with plain secret for return
	service := &domain.Services{
		ServiceName:  req.Name,
		Description:  &req.Description,
		Logo:         &logoPath,
		ClientId:     clientID,
		ClientSecret: clientSecret,
		RedirectUrl:  req.RedirectURL,
		IsActive:     true,
	}

	// Create a copy for database with hashed secret
	hashedService := *service
	hashedService.ClientSecret = string(hashedSecret)

	if err := t.serviceRepository.CreateClient(&hashedService); err != nil {
		return nil, err
	}


	service.ServiceId = hashedService.ServiceId
	return service, nil
}

func (t *serviceUsecase) DeleteClientById(serviceId uint64) error {
	_, err := t.serviceRepository.GetClientById(serviceId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.ErrNotFound
		}
		return pkgErrors.ErrFailedQuery
	}

	if err := t.serviceRepository.DeleteClientById(serviceId); err != nil {
		return pkgErrors.ErrFailedDelete
	}

	return nil
}

func (t *serviceUsecase) GetClients(p *pagination.Pagination) (*pagination.Pagination, error) {
	clients, err := t.serviceRepository.GetClients(p)
	if err != nil {
		return nil, pkgErrors.ErrFailedQuery
	}

	return clients, nil
}

func (t *serviceUsecase) UpdateClient(req request.UpdateClientReq, file multipart.File, header *multipart.FileHeader) error {
	service, err := t.serviceRepository.GetClientById(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.ErrNotFound
		}
		return pkgErrors.ErrFailedQuery
	}

	service.ServiceName = req.Name
	service.Description = &req.Description
	service.RedirectUrl = req.RedirectURL
	if req.IsActive != nil {
		service.IsActive = *req.IsActive
	}

	var oldLogoPath string
	if file != nil {
		if service.Logo != nil {
			oldLogoPath = *service.Logo
		}

		logoPath, err := image.ProcessImage(file, header, "public/images/services")
		if err != nil {
			return err
		}
		service.Logo = &logoPath
	}

	if err := t.serviceRepository.UpdateClient(service); err != nil {
		return pkgErrors.ErrFailedQuery
	}

	// Delete old logo after successful update
	if oldLogoPath != "" {
		// Remove leading slash to make path relative to CWD
		targetPath := oldLogoPath
		if len(targetPath) > 0 && (targetPath[0] == '/' || targetPath[0] == '\\') {
			targetPath = targetPath[1:]
		}
		
		// Normalize path for OS (handles Windows backslashes)
		targetPath = filepath.FromSlash(targetPath)
		
		// Delete file
		_ = os.Remove(targetPath)
	}

	return nil
}

func generateClientID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateClientSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
