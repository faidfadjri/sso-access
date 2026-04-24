package role_permission_test

import (
	"akastra-access/internal/app/usecases/role_permission"
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

func TestAssignPermission(t *testing.T) {
	mockRepo := new(mocks.RolePermissionRepositoryMock)
	usecase := role_permission.NewRolePermissionUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Create", mock.AnythingOfType("*domain.RolePermissions")).Return(nil).Once()

		rp, err := usecase.AssignPermission(1, 2)

		assert.NoError(t, err)
		assert.NotNil(t, rp)
		assert.Equal(t, uint64(1), rp.RoleId)
		assert.Equal(t, uint64(2), rp.PermissionId)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("Create", mock.AnythingOfType("*domain.RolePermissions")).Return(errors.New("db error")).Once()

		rp, err := usecase.AssignPermission(1, 2)

		assert.Error(t, err)
		assert.Nil(t, rp)
		assert.Equal(t, pkgErrors.ErrFailedQuery, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetList_RolePermission(t *testing.T) {
	mockRepo := new(mocks.RolePermissionRepositoryMock)
	usecase := role_permission.NewRolePermissionUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		expectedPagination := &pagination.Pagination{}
		mockRepo.On("GetList", mock.Anything).Return(expectedPagination, nil).Once()

		p := &pagination.Pagination{}
		res, err := usecase.GetList(p)

		assert.NoError(t, err)
		assert.Equal(t, expectedPagination, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetByKeys_RolePermission(t *testing.T) {
	mockRepo := new(mocks.RolePermissionRepositoryMock)
	usecase := role_permission.NewRolePermissionUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		expected := &domain.RolePermissions{RoleId: 1}
		mockRepo.On("GetByKeys", uint64(1), uint64(2)).Return(expected, nil).Once()

		res, err := usecase.GetByKeys(1, 2)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByKeys", uint64(1), uint64(2)).Return(nil, gorm.ErrRecordNotFound).Once()

		res, err := usecase.GetByKeys(1, 2)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, pkgErrors.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestRemovePermission(t *testing.T) {
	mockRepo := new(mocks.RolePermissionRepositoryMock)
	usecase := role_permission.NewRolePermissionUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetByKeys", uint64(1), uint64(2)).Return(&domain.RolePermissions{}, nil).Once()
		mockRepo.On("Delete", uint64(1), uint64(2)).Return(nil).Once()

		err := usecase.RemovePermission(1, 2)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByKeys", uint64(1), uint64(2)).Return(nil, gorm.ErrRecordNotFound).Once()

		err := usecase.RemovePermission(1, 2)

		assert.Error(t, err)
		assert.Equal(t, pkgErrors.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}
