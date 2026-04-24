package bootstrap

import (
	"akastra-access/internal/app/usecases"
	database "akastra-access/internal/infrastructure/databases"
	"akastra-access/internal/infrastructure/middlewares"
	"akastra-access/internal/infrastructure/repositories"
	handler "akastra-access/internal/interface/handlers"
	"log"

	"gorm.io/gorm"
)

type Dependencies struct {
	ServiceHandler         *handler.ServiceHandler
	OAuthHandler           *handler.OAuthHandler
	UserHandler            *handler.UserHandler
	RoleHandler            *handler.RoleHandler
	PermissionHandler      *handler.PermissionHandler
	UserAccessHandler      *handler.UserAccessHandler
	AssignedRolesHandler   *handler.AssignedRolesHandler
	RolePermissionHandler  *handler.RolePermissionHandler
	DB                     *gorm.DB
	SessionRepo            repositories.SessionRepository
	Middleware             *middlewares.Middleware
}

func InitDependencies() *Dependencies {

    db, err := database.ConnectDB()

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Redis
	redisClient := database.ConnectRedis()

	// Repositories
	serviceRepo := repositories.NewServiceRepository(db)
	oauthRepo := repositories.NewOAuthRepository(db)
	sessionRepo := repositories.NewSessionRepository(redisClient)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)
	userAccessRepo := repositories.NewUserAccessRepository(db)
	assignedRolesRepo := repositories.NewAssignedRolesRepository(db)
	rolePermissionRepo := repositories.NewRolePermissionRepository(db)

	// Usecases
	serviceUsecase := usecases.NewServiceUsecase(serviceRepo)
	oauthUsecase := usecases.NewOAuthUsecase(oauthRepo, serviceRepo, sessionRepo, assignedRolesRepo, userAccessRepo)
	userUsecase := usecases.NewUserUsecase(userRepo, userAccessRepo)
	roleUsecase := usecases.NewRoleUsecase(roleRepo)
	permissionUsecase := usecases.NewPermissionUsecase(permissionRepo)
	userAccessUsecase := usecases.NewUserAccessUsecase(db, userAccessRepo)
	assignedRolesUsecase := usecases.NewAssignedRolesUsecase(assignedRolesRepo)
	rolePermissionUsecase := usecases.NewRolePermissionUsecase(rolePermissionRepo)

	// Handlers
	serviceHandler := handler.NewServiceHandler(serviceUsecase)
	oauthHandler := handler.NewOAuthHandler(oauthUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	roleHandler := handler.NewRoleHandler(roleUsecase)
	permissionHandler := handler.NewPermissionHandler(permissionUsecase)
	userAccessHandler := handler.NewUserAccessHandler(userAccessUsecase)
	assignedRolesHandler := handler.NewAssignedRolesHandler(assignedRolesUsecase)
	rolePermissionHandler := handler.NewRolePermissionHandler(rolePermissionUsecase)

	// Middleware
	middleware := middlewares.NewMiddleware(sessionRepo, oauthRepo)

	return &Dependencies{
		ServiceHandler:         serviceHandler,
		OAuthHandler:           oauthHandler,
		UserHandler:            userHandler,
		RoleHandler:            roleHandler,
		PermissionHandler:      permissionHandler,
		UserAccessHandler:      userAccessHandler,
		AssignedRolesHandler:   assignedRolesHandler,
		RolePermissionHandler:  rolePermissionHandler,
		DB:                     db,
		SessionRepo:            sessionRepo,
		Middleware:             middleware,
	}
}
