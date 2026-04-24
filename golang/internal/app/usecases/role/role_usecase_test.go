package role_test

import (
	"akastra-access/internal/app/usecases/role"
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"
	"akastra-access/internal/mocks"
	pkgErrors "akastra-access/internal/pkg/errors"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateRole(t *testing.T) {
	mockRepo := new(mocks.RoleRepositoryMock)
	usecase := role.NewRoleUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Create", mock.AnythingOfType("*domain.ServiceRoles")).Return(nil).Once()

		role, err := usecase.CreateRole("Admin", 1)

		assert.NoError(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, "Admin", role.RoleName)
		assert.Equal(t, uint64(1), role.ServiceId)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("Create", mock.AnythingOfType("*domain.ServiceRoles")).Return(errors.New("db error")).Once()

		role, err := usecase.CreateRole("Admin", 1)

		assert.Error(t, err)
		assert.Nil(t, role)
		assert.Equal(t, pkgErrors.ErrFailedQuery, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetRoles(t *testing.T) {
	mockRepo := new(mocks.RoleRepositoryMock)
	usecase := role.NewRoleUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		expectedPagination := &pagination.Pagination{}
		mockRepo.On("GetRoles", mock.Anything).Return(expectedPagination, nil).Once()

		p := &pagination.Pagination{}
		res, err := usecase.GetRoles(p)

		assert.NoError(t, err)
		assert.Equal(t, expectedPagination, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetRoles", mock.Anything).Return(nil, errors.New("db error")).Once()

		p := &pagination.Pagination{}
		res, err := usecase.GetRoles(p)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, pkgErrors.ErrFailedQuery, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetRoleByID(t *testing.T) {
	mockRepo := new(mocks.RoleRepositoryMock)
	usecase := role.NewRoleUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		expectedRole := &domain.ServiceRoles{RoleName: "Admin"}
		mockRepo.On("GetRoleByID", uint64(1)).Return(expectedRole, nil).Once()

		res, err := usecase.GetRoleByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedRole, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetRoleByID", uint64(1)).Return(nil, gorm.ErrRecordNotFound).Once()

		res, err := usecase.GetRoleByID(1)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, pkgErrors.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateRole(t *testing.T) {
	mockRepo := new(mocks.RoleRepositoryMock)
	usecase := role.NewRoleUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		existingRole := &domain.ServiceRoles{RoleName: "Old"}
		mockRepo.On("GetRoleByID", uint64(1)).Return(existingRole, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*domain.ServiceRoles")).Return(nil).Once()

		err := usecase.UpdateRole(1, "New", 2)

		assert.NoError(t, err)
		assert.Equal(t, "New", existingRole.RoleName)
		assert.Equal(t, uint64(2), existingRole.ServiceId)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetRoleByID", uint64(1)).Return(nil, gorm.ErrRecordNotFound).Once()

		err := usecase.UpdateRole(1, "New", 2)

		assert.Error(t, err)
		assert.Equal(t, pkgErrors.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteRole(t *testing.T) {
	mockRepo := new(mocks.RoleRepositoryMock)
	usecase := role.NewRoleUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetRoleByID", uint64(1)).Return(&domain.ServiceRoles{}, nil).Once()
		mockRepo.On("Delete", uint64(1)).Return(nil).Once()

		err := usecase.DeleteRole(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetRoleByID", uint64(1)).Return(nil, gorm.ErrRecordNotFound).Once()

		err := usecase.DeleteRole(1)

		assert.Error(t, err)
		assert.Equal(t, pkgErrors.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}
