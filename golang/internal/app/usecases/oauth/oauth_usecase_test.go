package oauth_test

import (
	"akastra-access/internal/app/entities"
	"akastra-access/internal/app/usecases/oauth"
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/interface/http/request"
	"akastra-access/internal/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	mockOAuthRepo := new(mocks.OAuthRepositoryMock)
	mockServiceRepo := new(mocks.ServiceRepositoryMock)
	mockSessionRepo := new(mocks.SessionRepositoryMock)
	mockAssignedRolesRepo := new(mocks.AssignedRolesRepositoryMock)
	mockUserAccessRepo := new(mocks.UserAccessRepositoryMock)
	usecase := oauth.NewOAuthUsecase(mockOAuthRepo, mockServiceRepo, mockSessionRepo, mockAssignedRolesRepo, mockUserAccessRepo)

	t.Run("success", func(t *testing.T) {
		password := "password"
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := &domain.Users{UserId: 1, Password: string(hashed)}
		
		req := request.LoginRequest{EmailOrUsername: "test", Password: password}

		mockOAuthRepo.On("GetUserByEmailorUsername", "test").Return(user, nil).Once()
		mockSessionRepo.On("SetSession", mock.Anything, uint64(1), 24*time.Hour).Return(nil).Once()

		token, _, err := usecase.Login(req)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockOAuthRepo.AssertExpectations(t)
		mockSessionRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		user := &domain.Users{UserId: 1, Password: string(hashed)}
		
		req := request.LoginRequest{EmailOrUsername: "test", Password: "wrong_password"}

		mockOAuthRepo.On("GetUserByEmailorUsername", "test").Return(user, nil).Once()

		token, _, err := usecase.Login(req)

		assert.Error(t, err)
		assert.Empty(t, token)
		mockOAuthRepo.AssertExpectations(t)
	})
}

func TestUpdateAccount(t *testing.T) {
	mockOAuthRepo := new(mocks.OAuthRepositoryMock)
	mockServiceRepo := new(mocks.ServiceRepositoryMock)
	mockSessionRepo := new(mocks.SessionRepositoryMock)
	mockAssignedRolesRepo := new(mocks.AssignedRolesRepositoryMock)
	mockUserAccessRepo := new(mocks.UserAccessRepositoryMock)
	usecase := oauth.NewOAuthUsecase(mockOAuthRepo, mockServiceRepo, mockSessionRepo, mockAssignedRolesRepo, mockUserAccessRepo)

	t.Run("success update profile without photo", func(t *testing.T) {
		user := &domain.Users{UserId: 1, FullName: "Old Name", Email: "old@test.com", Username: "olduser"}
		
		req := request.UpdateAccountRequest{
			FullName: "New Name",
			Email:    "new@test.com",
			Username: "newuser",
			Phone:    "1234567890",
		}

		mockOAuthRepo.On("GetUserByID", uint64(1)).Return(user, nil).Once()
		mockOAuthRepo.On("UpdateUser", mock.AnythingOfType("*domain.Users")).Return(nil).Once()
		mockAssignedRolesRepo.On("GetByKeys", uint64(1), mock.Anything).Return(&entities.AssignedRoles{RoleName: "user"}, nil).Once()

		_, _, err := usecase.UpdateAccount(1, req, nil, nil)

		assert.NoError(t, err)
		assert.Equal(t, "New Name", user.FullName)
		assert.Equal(t, "new@test.com", user.Email)
		mockOAuthRepo.AssertExpectations(t)
		mockAssignedRolesRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		req := request.UpdateAccountRequest{}
		mockOAuthRepo.On("GetUserByID", uint64(1)).Return(nil, nil).Once()

		_, _, err := usecase.UpdateAccount(1, req, nil, nil)

		assert.Error(t, err)
		mockOAuthRepo.AssertExpectations(t)
	})
}
