package permission_test

import (
	"akastra-access/internal/app/usecases/permission"
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

func TestCreatePermission(t *testing.T) {
	mockRepo := new(mocks.PermissionRepositoryMock)
	usecase := permission.NewPermissionUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Create", mock.AnythingOfType("*domain.Permissions")).Return(nil).Once()

		perm, err := usecase.CreatePermission("read:users", "Read access")

		assert.NoError(t, err)
		assert.NotNil(t, perm)
		assert.Equal(t, "read:users", perm.PermissionKey)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("Create", mock.AnythingOfType("*domain.Permissions")).Return(errors.New("db error")).Once()

		perm, err := usecase.CreatePermission("read:users", "Read access")

		assert.Error(t, err)
		assert.Nil(t, perm)
		assert.Equal(t, pkgErrors.ErrFailedQuery, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetPermissions(t *testing.T) {
	mockRepo := new(mocks.PermissionRepositoryMock)
	usecase := permission.NewPermissionUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		expectedPagination := &pagination.Pagination{}
		mockRepo.On("GetPermissions", mock.Anything).Return(expectedPagination, nil).Once()

		p := &pagination.Pagination{}
		res, err := usecase.GetPermissions(p)

		assert.NoError(t, err)
		assert.Equal(t, expectedPagination, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetPermissionByID(t *testing.T) {
	mockRepo := new(mocks.PermissionRepositoryMock)
	usecase := permission.NewPermissionUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		expectedPerm := &domain.Permissions{PermissionKey: "read"}
		mockRepo.On("GetPermissionByID", uint64(1)).Return(expectedPerm, nil).Once()

		res, err := usecase.GetPermissionByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedPerm, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetPermissionByID", uint64(1)).Return(nil, gorm.ErrRecordNotFound).Once()

		res, err := usecase.GetPermissionByID(1)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, pkgErrors.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdatePermission(t *testing.T) {
	mockRepo := new(mocks.PermissionRepositoryMock)
	usecase := permission.NewPermissionUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		existingPerm := &domain.Permissions{PermissionKey: "old"}
		mockRepo.On("GetPermissionByID", uint64(1)).Return(existingPerm, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*domain.Permissions")).Return(nil).Once()

		err := usecase.UpdatePermission(1, "new", "desc")

		assert.NoError(t, err)
		assert.Equal(t, "new", existingPerm.PermissionKey)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeletePermission(t *testing.T) {
	mockRepo := new(mocks.PermissionRepositoryMock)
	usecase := permission.NewPermissionUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetPermissionByID", uint64(1)).Return(&domain.Permissions{}, nil).Once()
		mockRepo.On("Delete", uint64(1)).Return(nil).Once()

		err := usecase.DeletePermission(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
