package oauth

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/interface/http/request"
	"mime/multipart"
)

// OAuthUsecase defines the interface for OAuth use case
type OAuthUsecase interface {
	Login(req request.LoginRequest) (string, *domain.Users, error)
	RefreshToken(req request.RefreshTokenRequest) (string, string, error)
	Authorize(req request.AuthorizeRequest, sessionID string) (string, error)
	TokenExchange(req request.ExchangeAuthCodeRequest) (string, string, error)
	Logout(sessionID string) error
	UpdateAccount(userId uint64, req request.UpdateAccountRequest, file multipart.File, header *multipart.FileHeader) (string, string, error)
	ForgotPassword(req request.ForgotPasswordRequest) error
	ResetPassword(req request.ResetPasswordRequest) error
}
