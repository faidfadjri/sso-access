package oauth_test

import (
	"akastra-access/internal/app/entities"
	"akastra-access/internal/app/usecases/oauth"
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/interface/http/request"
	"akastra-access/internal/mocks"

	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestTokenExchange(t *testing.T) {
	mockOAuthRepo := new(mocks.OAuthRepositoryMock)
	mockServiceRepo := new(mocks.ServiceRepositoryMock)
	mockSessionRepo := new(mocks.SessionRepositoryMock)
	mockAssignedRolesRepo := new(mocks.AssignedRolesRepositoryMock)
	mockUserAccessRepo := new(mocks.UserAccessRepositoryMock)
	usecase := oauth.NewOAuthUsecase(mockOAuthRepo, mockServiceRepo, mockSessionRepo, mockAssignedRolesRepo, mockUserAccessRepo)

	t.Run("success", func(t *testing.T) {
		secret := "secret"
		hashedSecret, _ := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
		
		req := request.ExchangeAuthCodeRequest{
			GrantType:    "authorization_code",
			ClientID:     "client123",
			RedirectURI:  "http://localhost:3000/callback",
			Code:         "auth_code_123",
			ClientSecret: toPtr(secret),
		}

		// Mock responses
		service := &domain.Services{ServiceId: 1, ClientId: "client123", ClientSecret: string(hashedSecret), RedirectUrl: "http://localhost:3000/callback"} 
		
		// AuthSession
		sessionJSON := `{"client_id":"client123","user_id":1,"redirect_uri":"http://localhost:3000/callback"}`
		
		// UserWithService
		userWithService := &entities.UserWithService{UserId: 1, Email: "test@test.com"}
		
		// AssignedRole (IDP Role)
		idpRole := &entities.AssignedRoles{RoleName: "Admin"}

		mockServiceRepo.On("GetClientByClientId", "client123").Return(service, nil).Once()
		mockSessionRepo.On("GetAuthCode", "auth_code_123").Return(sessionJSON, nil).Once()
		
		mockOAuthRepo.On("GetUserByIDWithService", uint64(1), uint64(1)).Return(userWithService, nil).Once()
		
		// IDP ID is "1" by default in config.
		mockAssignedRolesRepo.On("GetByKeys", uint64(1), uint64(1)).Return(idpRole, nil).Once()
		
		mockUserAccessRepo.On("IsUserAccessExist", uint64(1), uint64(1)).Return(true, nil).Once()

		mockSessionRepo.On("DeleteAuthCode", "auth_code_123").Return(nil).Once()

		accessToken, refreshToken, err := usecase.TokenExchange(req)

		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
		mockOAuthRepo.AssertExpectations(t)
		mockServiceRepo.AssertExpectations(t)
		mockSessionRepo.AssertExpectations(t)
		mockAssignedRolesRepo.AssertExpectations(t)
		mockUserAccessRepo.AssertExpectations(t)
	})
}

func toPtr(s string) *string {
	return &s
}
