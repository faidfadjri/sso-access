package service

import (
	"akastra-access/internal/app/config"
	"akastra-access/internal/infrastructure/databases/domain"

	"akastra-access/internal/infrastructure/utils/pagination"

	"gorm.io/gorm"
)

// ServiceRepository handles data access
type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) CreateClient(req *domain.Services) error {
	return r.db.Create(req).Error
}

func (r *serviceRepository) GetClientById(serviceId uint64) (*domain.Services, error) {
	var service domain.Services
	err := r.db.Where("service_id = ?", serviceId).First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepository) GetClientByIds(serviceIds []uint64) ([]*domain.Services, error) {
	var services []*domain.Services
	err := r.db.Where("service_id IN ?", serviceIds).Find(&services).Error
	if err != nil {
		return nil, err
	}
	return services, nil
}

func (r *serviceRepository) GetClientByClientId(clientId string) (*domain.Services, error) {
	var service domain.Services
	err := r.db.Where("client_id = ?", clientId).First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepository) DeleteClientById(serviceId uint64) error {
	return r.db.Where("service_id = ?", serviceId).Delete(&domain.Services{}).Error
}

func (r *serviceRepository) GetClients(p *pagination.Pagination) (*pagination.Pagination, error) {
	var services []domain.Services

	db := r.db.Model(&domain.Services{})

	if p.Search != "" {
		db = db.Where("service_name LIKE ?", "%"+p.Search+"%")
	}

	frontendURL := config.Load().FrontendURL
	db = db.Where("redirect_url NOT LIKE ?", "%"+frontendURL+"%")
	err := db.Scopes(p.Paginate(&services, db)).Find(&services).Error
	if err != nil {
		return nil, err
	}

	p.Rows = services
	p.Rows = services
	return p, nil
}

func (r *serviceRepository) UpdateClient(req *domain.Services) error {
	return r.db.Save(req).Error
}
