package errors

import "errors"

var (
	ErrInvalidFormat 		= errors.New("invalid file format, only JPEG and PNG are allowed")
	ErrNotFound      		= errors.New("resource not found")
	ErrInternal      		= errors.New("internal server error")
	ErrBadRequest    		= errors.New("bad request")
	ErrUnauthorized  		= errors.New("unauthorized")
	ErrFailedQuery 		    = errors.New("failed to query database")
	ErrFailedDelete 		= errors.New("failed to delete resource")
	
	// Oauth
	ErrInvalidBearerToken   = errors.New("invalid bearer token")
	ErrInvalidClientID 		= errors.New("invalid client id")
	ErrInvalidRedirectURI 	= errors.New("invalid redirect uri")
	ErrInvalidResponseType 	= errors.New("invalid response type")
	ErrInvalidScope 		= errors.New("invalid scope")
	ErrInvalidAuthCode 		= errors.New("invalid or expired authorization code")
	ErrSessionNotFound 		= errors.New("session not found")
	ErrSessionExpired  		= errors.New("session expired")
	ErrInvalidClientSecret  = errors.New("invalid client secret")
	ErrInvalidCodeVerifier  = errors.New("invalid code verifier")


	// User Access
	ErrInvalidServiceIds = errors.New("invalid service ids")
	ErrInvalidUserId = errors.New("invalid user id")
	ErrInvalidRoleId = errors.New("invalid role id")
	ErrUserNotHaveAccess = errors.New("user not have access")
)
