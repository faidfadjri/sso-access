package assigned_roles

import (
	"akastra-access/internal/app/entities"
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

func TestAssignRole(t *testing.T) {
	mockRepo := new(mocks.AssignedRolesRepositoryMock)
	usecase := NewAssignedRolesUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetServiceIdByRoleId", uint64(2)).Return(uint64(10), nil).Once()
		mockRepo.On("Create", mock.AnythingOfType("*domain.UserServiceRole")).Return(nil).Once()

		usrList, err := usecase.AssignRole([]uint64{1}, 2)

		assert.NoError(t, err)
		assert.NotNil(t, usrList)
		assert.Len(t, usrList, 1)
		assert.Equal(t, uint64(1), usrList[0].UserId)
		assert.Equal(t, uint64(2), usrList[0].RoleId)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetServiceIdByRoleId", uint64(2)).Return(uint64(10), nil).Once()
		mockRepo.On("Create", mock.AnythingOfType("*domain.UserServiceRole")).Return(errors.New("db error")).Once()

		usrList, err := usecase.AssignRole([]uint64{1}, 2)

		assert.Error(t, err)
		assert.Nil(t, usrList)
		assert.Equal(t, pkgErrors.ErrFailedQuery, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("duplicate_ignored", func(t *testing.T) {
		mockRepo.On("GetServiceIdByRoleId", uint64(2)).Return(uint64(10), nil).Once()
		// Mock Create for first user (duplicate)
		mockRepo.On("Create", mock.MatchedBy(func(usr *domain.UserServiceRole) bool {
			return usr.UserId == 1
		})).Return(errors.New("Duplicate entry")).Once()
		// Mock Create for second user (success)
		mockRepo.On("Create", mock.MatchedBy(func(usr *domain.UserServiceRole) bool {
			return usr.UserId == 3
		})).Return(nil).Once()

		usrList, err := usecase.AssignRole([]uint64{1, 3}, 2)

		assert.NoError(t, err)
		assert.NotNil(t, usrList)
		assert.Len(t, usrList, 1) // Only 1 successfully created
		assert.Equal(t, uint64(3), usrList[0].UserId)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetList(t *testing.T) {
	mockRepo := new(mocks.AssignedRolesRepositoryMock)
	usecase := NewAssignedRolesUsecase(mockRepo)

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

func TestGetByKeys(t *testing.T) {
	mockRepo := new(mocks.AssignedRolesRepositoryMock)
	usecase := NewAssignedRolesUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		expected := &entities.AssignedRoles{UserId: 1}
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

func TestRemoveRole(t *testing.T) {
	mockRepo := new(mocks.AssignedRolesRepositoryMock)
	usecase := NewAssignedRolesUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetByKeys", uint64(1), uint64(2)).Return(&entities.AssignedRoles{}, nil).Once()
		mockRepo.On("Delete", uint64(1), uint64(2)).Return(nil).Once()

		err := usecase.RemoveRole(1, 2)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByKeys", uint64(1), uint64(2)).Return(nil, gorm.ErrRecordNotFound).Once()

		err := usecase.RemoveRole(1, 2)

		assert.Error(t, err)
		assert.Equal(t, pkgErrors.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}
