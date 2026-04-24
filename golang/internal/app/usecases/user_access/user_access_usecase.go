package user_access

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/repositories"
	repoUserAccess "akastra-access/internal/infrastructure/repositories/user_access"
	"akastra-access/internal/infrastructure/utils/pagination"
	pkgErrors "akastra-access/internal/pkg/errors"
	"errors"
	"log"

	"gorm.io/gorm"
)

type userAccessUsecase struct {
	db *gorm.DB
	accessRepo repoUserAccess.UserAccessRepository
}

func NewUserAccessUsecase(db *gorm.DB, repo repoUserAccess.UserAccessRepository) UserAccessUsecase {
	return &userAccessUsecase{
		db: db,
		accessRepo: repo,
	}
}

func (u *userAccessUsecase) CreateUserAccess(
	userID uint64,
	serviceIDs []uint64,
	status string,
) ([]*domain.UserAccess, error) {

	var accesses []*domain.UserAccess

	err := u.db.Transaction(func(tx *gorm.DB) error {
		userAccessRepo := repositories.NewUserAccessRepository(tx)
		serviceRepo := repositories.NewServiceRepository(tx)
		userRepo := repositories.NewUserRepository(tx)

		// validate services
		services, err := serviceRepo.GetClientByIds(serviceIDs)
		if err != nil {
			return pkgErrors.ErrFailedQuery
		}

		if len(services) != len(serviceIDs) {
			return pkgErrors.ErrInvalidServiceIds
		}

		// validate user
		if _, err := userRepo.GetUserByID(userID); err != nil {
			return pkgErrors.ErrInvalidUserId
		}

		// get existing accesses
		accessList, err := userAccessRepo.GetUserAccessByUserId(userID)
		if err != nil {
			return pkgErrors.ErrFailedQuery
		}

		// build set of existing service IDs
		existing := make(map[uint64]struct{})
		for _, access := range accessList {
			existing[access.ServiceId] = struct{}{}
		}

		// filter only new service IDs
		var newServiceIDs []uint64
		for _, serviceID := range serviceIDs {
			if _, found := existing[serviceID]; !found {
				newServiceIDs = append(newServiceIDs, serviceID)
			}
		}

		// create accesses
		for _, serviceID := range newServiceIDs {
			accessData := &domain.UserAccess{
				UserId:    userID,
				ServiceId: serviceID,
				Status:    status,
			}

			if err := userAccessRepo.Create(accessData); err != nil {
				return pkgErrors.ErrFailedQuery
			}

			accesses = append(accesses, accessData)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return accesses, nil
}


func (u *userAccessUsecase) GetUserAccesses(p *pagination.Pagination) (*pagination.Pagination, error) {
	accesses, err := u.accessRepo.GetUserAccesses(p)
	log.Println(accesses)
	if err != nil {
		return nil, pkgErrors.ErrFailedQuery
	}
	return accesses, nil
}

func (u *userAccessUsecase) GetUserAccessByID(id uint64) (*domain.UserAccess, error) {
	access, err := u.accessRepo.GetUserAccessByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, pkgErrors.ErrFailedQuery
	}
	return access, nil
}

func (u *userAccessUsecase) UpdateUserAccess(userID uint64, serviceIDs []uint64, status string) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		userAccessRepo := repositories.NewUserAccessRepository(tx)
		serviceRepo := repositories.NewServiceRepository(tx)
		userRepo := repositories.NewUserRepository(tx)

		// validate services
		services, err := serviceRepo.GetClientByIds(serviceIDs)
		if err != nil {
			return pkgErrors.ErrFailedQuery
		}

		if len(services) != len(serviceIDs) {
			return pkgErrors.ErrInvalidServiceIds
		}

		// validate user
		if _, err := userRepo.GetUserByID(userID); err != nil {
			return pkgErrors.ErrInvalidUserId
		}

		// get existing accesses
		existingAccesses, err := userAccessRepo.GetUserAccessByUserId(userID)
		if err != nil {
			return pkgErrors.ErrFailedQuery
		}

		// map existing accesses by service ID
		existingMap := make(map[uint64]*domain.UserAccess)
		for i := range existingAccesses {
			existingMap[existingAccesses[i].ServiceId] = &existingAccesses[i]
		}

		// map new service IDs
		newMap := make(map[uint64]struct{})
		for _, id := range serviceIDs {
			newMap[id] = struct{}{}
		}

		// Identify to Add or Update
		for _, id := range serviceIDs {
			if _, exists := existingMap[id]; !exists {
				// Create new access
				newAccess := &domain.UserAccess{
					UserId:    userID,
					ServiceId: id,
					Status:    status,
				}
				if err := userAccessRepo.Create(newAccess); err != nil {
					return pkgErrors.ErrFailedQuery
				}
			} else {
				// Update status if exists
				access := existingMap[id]
				if access.Status != status {
					access.Status = status
					if err := userAccessRepo.Update(access); err != nil {
						return pkgErrors.ErrFailedQuery
					}
				}
			}
		}

		// Identify to Delete
		for _, access := range existingAccesses {
			if _, keep := newMap[access.ServiceId]; !keep {
				if err := userAccessRepo.Delete(access.AccessId); err != nil {
					return pkgErrors.ErrFailedDelete
				}
			}
		}

		return nil
	})
}

func (u *userAccessUsecase) DeleteUserAccess(id uint64) error {
	if _, err := u.accessRepo.GetUserAccessByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.ErrNotFound
		}
		return pkgErrors.ErrFailedQuery
	}

	if err := u.accessRepo.Delete(id); err != nil {
		return pkgErrors.ErrFailedDelete
	}

	return nil
}
