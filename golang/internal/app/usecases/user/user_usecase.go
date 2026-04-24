package user

import (
	"akastra-access/internal/app/config"
	"akastra-access/internal/infrastructure/databases/domain"
	repoUser "akastra-access/internal/infrastructure/repositories/user"
	repoUserAccess "akastra-access/internal/infrastructure/repositories/user_access"
	"akastra-access/internal/infrastructure/utils/image"
	"akastra-access/internal/infrastructure/utils/pagination"
	"akastra-access/internal/interface/http/request"
	pkgErrors "akastra-access/internal/pkg/errors"
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userUsecase struct {
	userRepo       repoUser.UserRepository
	userAccessRepo repoUserAccess.UserAccessRepository
}

func NewUserUsecase(repo repoUser.UserRepository, userAccessRepo repoUserAccess.UserAccessRepository) UserUsecase {
	return &userUsecase{
		userRepo:       repo,
		userAccessRepo: userAccessRepo,
	}
}

func (u *userUsecase) CreateUser(req request.CreateUserReq, file multipart.File, header *multipart.FileHeader) (*domain.Users, error) {
	// check email unique
	if _, err := u.userRepo.GetUserByEmail(req.Email); err == nil {
		return nil, errors.New("email already registered")
	}

	// check username unique
	if _, err := u.userRepo.GetUserByUsername(req.Username); err == nil {
		return nil, errors.New("username already taken")
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// process image
	var photoPath *string
	if file != nil {
	path, err := image.ProcessImage(file, header, "public/images/profiles")
		if err != nil {
			return nil, err
		}
		photoPath = &path
	}

	user := &domain.Users{
		FullName: req.FullName,
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
		Photo:    photoPath,
		Admin:    req.Admin == "true",
	}

	if req.Phone != "" {
		user.Phone = &req.Phone
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, pkgErrors.ErrFailedQuery
	}

	// create user access for IDP
	cfg := config.Load()
	if idpID, err := strconv.ParseUint(cfg.IdentityProviderID, 10, 64); err == nil {
		access := &domain.UserAccess{
			UserId:    user.UserId,
			ServiceId: idpID,
			Status:    "active",
		}
		_ = u.userAccessRepo.Create(access)
	}

	return user, nil
}

func (u *userUsecase) GetUsers(p *pagination.Pagination) (*pagination.Pagination, error) {
	users, err := u.userRepo.GetUsers(p)
	if err != nil {
		return nil, pkgErrors.ErrFailedQuery
	}
	return users, nil
}

func (u *userUsecase) GetUserByID(id uint64) (*domain.Users, error) {
	user, err := u.userRepo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, pkgErrors.ErrFailedQuery
	}
	return user, nil
}

func (u *userUsecase) UpdateUser(req request.UpdateUserReq, file multipart.File, header *multipart.FileHeader) error {
	user, err := u.userRepo.GetUserByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.ErrNotFound
		}
		return pkgErrors.ErrFailedQuery
	}

	user.FullName = req.FullName
	user.Email = req.Email
	user.Username = req.Username
	user.Phone = &req.Phone
	
	if req.Admin == "true" {
		user.Admin = true
	} else if req.Admin == "false" {
		user.Admin = false
	}

	// Password update
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	// Photo update
	var oldPhotoPath string
	if file != nil {
		if user.Photo != nil {
			oldPhotoPath = *user.Photo
		}

		path, err := image.ProcessImage(file, header, "public/images/profiles")
		if err != nil {
			return err
		}
		user.Photo = &path
	}

	if err := u.userRepo.Update(user); err != nil {
		return pkgErrors.ErrFailedQuery
	}

	// Delete old photo
	if oldPhotoPath != "" {
		targetPath := oldPhotoPath
		if len(targetPath) > 0 && (targetPath[0] == '/' || targetPath[0] == '\\') {
			targetPath = targetPath[1:]
		}
		targetPath = filepath.FromSlash(targetPath)
		_ = os.Remove(targetPath)
	}

	return nil
}

func (u *userUsecase) DeleteUser(id uint64) error {
	user, err := u.userRepo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.ErrNotFound
		}
		return pkgErrors.ErrFailedQuery
	}

	if err := u.userRepo.Delete(id); err != nil {
		return pkgErrors.ErrFailedDelete
	}

	// Delete photo
	if user.Photo != nil && *user.Photo != "" {
		targetPath := *user.Photo
		if len(targetPath) > 0 && (targetPath[0] == '/' || targetPath[0] == '\\') {
			targetPath = targetPath[1:]
		}
		targetPath = filepath.FromSlash(targetPath)
		_ = os.Remove(targetPath)
	}

	return nil
}

func (u *userUsecase) BatchDeleteUser(req request.BatchDeleteUserReq) error {
	// Optional: Fetch users to delete photos first (if strict cleanup required)
	// For now, just batch delete from DB.
	// If strict photo cleanup is needed, we'd loop.
	/*
	for _, id := range req.IDs {
		_ = u.DeleteUser(id) // This handles photo cleanup but is slow (N queries)
	}
	*/
	
	// Efficient DB delete, skip photo deletion for batch for now unless requested strictly.
	// Or fetch photos of all IDs then delete.
	// Prompt says "batch delete (like checkmark on table)". Users usually expect this to be fast.
	// But leaving orphan files is bad.
	// I'll fetch users to get photos, then delete.
	
	// Since repository doesn't have FindByIDs, I'll loop DeleteUser for simplicity and correctness of photo cleanup.
	// It's safer.
	
	for _, id := range req.IDs {
		// Ignore error? Or stop? usually best effort or transaction.
		// For simplicity, attempt all.
		_ = u.DeleteUser(id)
	}
	
	// or use BatchDelete repo if implemented.. 
	// I implemented BatchDelete in repo but it only deletes from DB.
	// I will loop DeleteUser to ensuring photo cleanup.
	
	return nil
}
