package repositories

import (
	"akastra-access/internal/infrastructure/repositories/assigned_roles"
	"akastra-access/internal/infrastructure/repositories/oauth"
	"akastra-access/internal/infrastructure/repositories/permission"
	"akastra-access/internal/infrastructure/repositories/role"
	"akastra-access/internal/infrastructure/repositories/role_permission"
	"akastra-access/internal/infrastructure/repositories/service"
	sessionrepository "akastra-access/internal/infrastructure/repositories/session"
	"akastra-access/internal/infrastructure/repositories/user"
	"akastra-access/internal/infrastructure/repositories/user_access"
)

type (
	OAuthRepository            = oauth.OAuthRepository
	ServiceRepository          = service.ServiceRepository
	SessionRepository          = sessionrepository.SessionRepository
	UserRepository             = user.UserRepository
	RoleRepository             = role.RoleRepository
	PermissionRepository       = permission.PermissionRepository
	UserAccessRepository       = user_access.UserAccessRepository
	AssignedRolesRepository    = assigned_roles.AssignedRolesRepository
	RolePermissionRepository   = role_permission.RolePermissionRepository
)

var (
	NewOAuthRepository            = oauth.NewOAuthRepository
	NewServiceRepository          = service.NewServiceRepository
	NewSessionRepository          = sessionrepository.NewSessionRepository
	NewUserRepository             = user.NewUserRepository
	NewRoleRepository             = role.NewRoleRepository
	NewPermissionRepository       = permission.NewPermissionRepository
	NewUserAccessRepository       = user_access.NewUserAccessRepository
	NewAssignedRolesRepository    = assigned_roles.NewAssignedRolesRepository
	NewRolePermissionRepository   = role_permission.NewRolePermissionRepository
)
