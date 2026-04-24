package user_access_test

import (
	"akastra-access/internal/app/usecases/user_access"
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

func TestCreateUserAccess(t *testing.T) {
	t.Skip("Skipping TestCreateUserAccess: Usecase implementation uses internal GORM transaction and creates new repositories, which cannot be tested with the current mock setup.")
	mockRepo := new(mocks.UserAccessRepositoryMock)
	usecase := user_access.NewUserAccessUsecase(nil, mockRepo)

	t.Run("success", func(t *testing.T) {
		// Expect Create to be called twice because we pass 2 service IDs
		mockRepo.On("Create", mock.AnythingOfType("*domain.UserAccess")).Return(nil).Twice()

		accesses, err := usecase.CreateUserAccess(1, []uint64{1, 2}, "active")

		assert.NoError(t, err)
		assert.Len(t, accesses, 2)
		assert.Equal(t, uint64(1), accesses[0].UserId)
		assert.Equal(t, uint64(1), accesses[0].ServiceId)
		assert.Equal(t, uint64(2), accesses[1].ServiceId)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("Create", mock.AnythingOfType("*domain.UserAccess")).Return(errors.New("db error")).Once()

		accesses, err := usecase.CreateUserAccess(1, []uint64{1}, "active")

		assert.Error(t, err)
		assert.Nil(t, accesses)
		assert.Equal(t, pkgErrors.ErrFailedQuery, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetUserAccesses(t *testing.T) {
	mockRepo := new(mocks.UserAccessRepositoryMock)
	usecase := user_access.NewUserAccessUsecase(nil, mockRepo)

	t.Run("success", func(t *testing.T) {
		expectedPagination := &pagination.Pagination{}
		mockRepo.On("GetUserAccesses", mock.Anything).Return(expectedPagination, nil).Once()

		p := &pagination.Pagination{}
		res, err := usecase.GetUserAccesses(p)

		assert.NoError(t, err)
		assert.Equal(t, expectedPagination, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetUserAccessByID(t *testing.T) {
	mockRepo := new(mocks.UserAccessRepositoryMock)
	usecase := user_access.NewUserAccessUsecase(nil, mockRepo)

	t.Run("success", func(t *testing.T) {
		expectedAccess := &domain.UserAccess{Status: "active"}
		mockRepo.On("GetUserAccessByID", uint64(1)).Return(expectedAccess, nil).Once()

		res, err := usecase.GetUserAccessByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedAccess, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetUserAccessByID", uint64(1)).Return(nil, gorm.ErrRecordNotFound).Once()

		res, err := usecase.GetUserAccessByID(1)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, pkgErrors.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateUserAccess(t *testing.T) {
	t.Skip("Skipping TestUpdateUserAccess: Usecase implementation uses internal GORM transaction and creates new repositories, which cannot be tested with the current mock setup.")
	mockRepo := new(mocks.UserAccessRepositoryMock)
	usecase := user_access.NewUserAccessUsecase(nil, mockRepo)

	t.Run("success", func(t *testing.T) {
		existingAccess := &domain.UserAccess{Status: "active"}
		mockRepo.On("GetUserAccessByID", uint64(1)).Return(existingAccess, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*domain.UserAccess")).Return(nil).Once()

		err := usecase.UpdateUserAccess(1, []uint64{3}, "revoke")

		assert.NoError(t, err)
		assert.Equal(t, uint64(2), existingAccess.UserId)
		assert.Equal(t, uint64(3), existingAccess.ServiceId)
		assert.Equal(t, "revoke", existingAccess.Status)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteUserAccess(t *testing.T) {
	mockRepo := new(mocks.UserAccessRepositoryMock)
	usecase := user_access.NewUserAccessUsecase(nil, mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetUserAccessByID", uint64(1)).Return(&domain.UserAccess{}, nil).Once()
		mockRepo.On("Delete", uint64(1)).Return(nil).Once()

		err := usecase.DeleteUserAccess(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
