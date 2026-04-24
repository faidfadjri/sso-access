package user

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"
	"akastra-access/internal/interface/http/request"
	"mime/multipart"
)

type UserUsecase interface {
	CreateUser(req request.CreateUserReq, file multipart.File, header *multipart.FileHeader) (*domain.Users, error)
	GetUsers(p *pagination.Pagination) (*pagination.Pagination, error)
	GetUserByID(id uint64) (*domain.Users, error)
	UpdateUser(req request.UpdateUserReq, file multipart.File, header *multipart.FileHeader) error
	DeleteUser(id uint64) error
	BatchDeleteUser(req request.BatchDeleteUserReq) error
}
