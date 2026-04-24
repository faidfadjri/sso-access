package usecases

import (
	"akastra-access/internal/app/usecases/assigned_roles"
	"akastra-access/internal/app/usecases/oauth"
	"akastra-access/internal/app/usecases/permission"
	"akastra-access/internal/app/usecases/role"
	"akastra-access/internal/app/usecases/role_permission"
	"akastra-access/internal/app/usecases/service"
	"akastra-access/internal/app/usecases/user"
	"akastra-access/internal/app/usecases/user_access"
)

type (
	OAuthUsecase            = oauth.OAuthUsecase
	ServiceUsecase          = service.ServiceUsecase
	UserUsecase             = user.UserUsecase
	RoleUsecase             = role.RoleUsecase
	PermissionUsecase       = permission.PermissionUsecase
	UserAccessUsecase       = user_access.UserAccessUsecase
	AssignedRolesUsecase    = assigned_roles.AssignedRolesUsecase
	RolePermissionUsecase   = role_permission.RolePermissionUsecase
)

var (
	NewOAuthUsecase            = oauth.NewOAuthUsecase
	NewServiceUsecase          = service.NewServiceUsecase
	NewUserUsecase             = user.NewUserUsecase
	NewRoleUsecase             = role.NewRoleUsecase
	NewPermissionUsecase       = permission.NewPermissionUsecase
	NewUserAccessUsecase       = user_access.NewUserAccessUsecase
	NewAssignedRolesUsecase = assigned_roles.NewAssignedRolesUsecase
	NewRolePermissionUsecase   = role_permission.NewRolePermissionUsecase
)
