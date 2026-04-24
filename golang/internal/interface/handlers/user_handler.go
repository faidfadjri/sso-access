package handler

import (
	"akastra-access/internal/app/usecases/user"
	"akastra-access/internal/infrastructure/utils/pagination"
	"akastra-access/internal/interface/http/request"
	"akastra-access/internal/interface/http/response"
	pkgErrors "akastra-access/internal/pkg/errors"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	usecase user.UserUsecase
}

func NewUserHandler(u user.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: u,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.BadRequest(w, "Invalid multipart form", err)
		return
	}

	file, header, _ := r.FormFile("photo")
	if file != nil {
		defer file.Close()
	}

	req := request.CreateUserReq{
		FullName: r.FormValue("full_name"),
		Email:    r.FormValue("email"),
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Phone:    r.FormValue("phone"),
		Admin:    r.FormValue("admin"),
	}

	if err := req.Validate(); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	user, err := h.usecase.CreateUser(req, file, header)
	if err != nil {
		if err.Error() == "email already registered" || err.Error() == "username already taken" {
			response.BadRequest(w, err.Error(), nil)
			return
		}
		log.Println(err.Error())
		response.InternalServerError(w, "Failed to create user", nil)
		return
	}

	response.Success(w, "User created successfully", user)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	pagination := pagination.GeneratePaginationFromRequest(r)

	users, err := h.usecase.GetUsers(pagination)
	if err != nil {
		response.InternalServerError(w, "Failed to get users", err)
		return
	}

	response.Success(w, "Users retrieved successfully", users)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := request.GetURLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid user ID", err)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.BadRequest(w, "Invalid multipart form", err)
		return
	}

	file, header, _ := r.FormFile("photo")
	if file != nil {
		defer file.Close()
	}

	req := request.UpdateUserReq{
		ID:              id,
		FullName:        r.FormValue("full_name"),
		Email:           r.FormValue("email"),
		Username:        r.FormValue("username"),
		Password:        r.FormValue("password"),
		PasswordConfirm: r.FormValue("password_confirmation"),
		Phone:           r.FormValue("phone"),
		Admin:           r.FormValue("admin"),
	}

	if err := req.Validate(); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	if err := h.usecase.UpdateUser(req, file, header); err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			response.NotFound(w, "User not found", err)
			return
		}
		log.Println(err.Error())
		response.InternalServerError(w, "Failed to update user", nil)
		return
	}

	response.Success(w, "User updated successfully", nil)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := request.GetURLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid user ID", err)
		return
	}

	if err := h.usecase.DeleteUser(id); err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			response.NotFound(w, "User not found", err)
			return
		}
		response.InternalServerError(w, "Failed to delete user", err)
		return
	}

	response.Success(w, "User deleted successfully", nil)
}

func (h *UserHandler) BatchDeleteUser(w http.ResponseWriter, r *http.Request) {
	var req request.BatchDeleteUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", err)
		return
	}

	if len(req.IDs) == 0 {
		response.BadRequest(w, "No user IDs provided", nil)
		return
	}

	if err := h.usecase.BatchDeleteUser(req); err != nil {
		response.InternalServerError(w, "Failed to batch delete users", err)
		return
	}

	response.Success(w, "Users deleted successfully", nil)
}
