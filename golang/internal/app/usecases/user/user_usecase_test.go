package user_test

import (
	"akastra-access/internal/app/usecases/user"
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"
	"akastra-access/internal/interface/http/request"
	"akastra-access/internal/mocks"
	pkgErrors "akastra-access/internal/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	mockAccessRepo := new(mocks.UserAccessRepositoryMock)
	usecase := user.NewUserUsecase(mockRepo, mockAccessRepo)

	t.Run("success", func(t *testing.T) {
		req := request.CreateUserReq{
			Email:    "test@example.com",
			Username: "testuser",
			Password: "password123",
			FullName: "Test User",
			Admin:    "false",
		}
		
		mockRepo.On("GetUserByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("GetUserByUsername", req.Username).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("Create", mock.AnythingOfType("*domain.Users")).Return(nil).Once()
		mockAccessRepo.On("Create", mock.AnythingOfType("*domain.UserAccess")).Return(nil).Maybe()

		res, err := usecase.CreateUser(req, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, req.Email, res.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("email taken", func(t *testing.T) {
		req := request.CreateUserReq{Email: "taken@example.com"}
		mockRepo.On("GetUserByEmail", req.Email).Return(&domain.Users{}, nil).Once()

		res, err := usecase.CreateUser(req, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "email already registered", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestGetUsers(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	mockAccessRepo := new(mocks.UserAccessRepositoryMock)
	usecase := user.NewUserUsecase(mockRepo, mockAccessRepo)

	t.Run("success", func(t *testing.T) {
		expectedPagination := &pagination.Pagination{}
		mockRepo.On("GetUsers", mock.Anything).Return(expectedPagination, nil).Once()

		p := &pagination.Pagination{}
		res, err := usecase.GetUsers(p)

		assert.NoError(t, err)
		assert.Equal(t, expectedPagination, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	mockAccessRepo := new(mocks.UserAccessRepositoryMock)
	usecase := user.NewUserUsecase(mockRepo, mockAccessRepo)

	t.Run("success", func(t *testing.T) {
		expectedUser := &domain.Users{Email: "test@example.com"}
		mockRepo.On("GetUserByID", uint64(1)).Return(expectedUser, nil).Once()

		res, err := usecase.GetUserByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetUserByID", uint64(1)).Return(nil, gorm.ErrRecordNotFound).Once()

		res, err := usecase.GetUserByID(1)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, pkgErrors.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	mockAccessRepo := new(mocks.UserAccessRepositoryMock)
	usecase := user.NewUserUsecase(mockRepo, mockAccessRepo)

	t.Run("success", func(t *testing.T) {
		existingUser := &domain.Users{Email: "old@example.com"}
		req := request.UpdateUserReq{ID: 1, Email: "new@example.com"}
		
		mockRepo.On("GetUserByID", uint64(1)).Return(existingUser, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*domain.Users")).Return(nil).Once()

		err := usecase.UpdateUser(req, nil, nil)

		assert.NoError(t, err)
		assert.Equal(t, "new@example.com", existingUser.Email)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	mockAccessRepo := new(mocks.UserAccessRepositoryMock)
	usecase := user.NewUserUsecase(mockRepo, mockAccessRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetUserByID", uint64(1)).Return(&domain.Users{}, nil).Once()
		mockRepo.On("Delete", uint64(1)).Return(nil).Once()

		err := usecase.DeleteUser(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
