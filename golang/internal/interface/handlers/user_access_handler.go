package handler

import (
	"akastra-access/internal/app/usecases/user_access"
	"akastra-access/internal/infrastructure/utils/pagination"
	"akastra-access/internal/interface/http/request"
	"akastra-access/internal/interface/http/response"
	pkgErrors "akastra-access/internal/pkg/errors"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type UserAccessHandler struct {
	usecase user_access.UserAccessUsecase
}

func NewUserAccessHandler(u user_access.UserAccessUsecase) *UserAccessHandler {
	return &UserAccessHandler{
		usecase: u,
	}
}

func (h *UserAccessHandler) CreateUserAccess(w http.ResponseWriter, r *http.Request) {
	var req request.CreateUserAccessReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", err)
		return
	}

	if err := req.Validate(); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	res, err := h.usecase.CreateUserAccess(req.UserId, req.ServiceIds, req.Status)
	if err != nil {
		response.InternalServerError(w, "Failed to create user access", err)
		return
	}

	response.Success(w, "User access created successfully", res)
}

func (h *UserAccessHandler) GetUserAccesses(w http.ResponseWriter, r *http.Request) {
	pagination := pagination.GeneratePaginationFromRequest(r)

	accesses, err := h.usecase.GetUserAccesses(pagination)
	if err != nil {
		response.InternalServerError(w, "Failed to get user accesses", err)
		return
	}

	response.Success(w, "User accesses retrieved successfully", accesses)
}

func (h *UserAccessHandler) GetUserAccessByID(w http.ResponseWriter, r *http.Request) {
	idStr := request.GetURLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid access ID", err)
		return
	}

	res, err := h.usecase.GetUserAccessByID(id)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			response.NotFound(w, "User access not found", err)
			return
		}
		response.InternalServerError(w, "Failed to get user access", err)
		return
	}

	response.Success(w, "User access retrieved successfully", res)
}

func (h *UserAccessHandler) UpdateUserAccess(w http.ResponseWriter, r *http.Request) {

	var req request.UpdateUserAccessReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", err)
		return
	}

	if err := req.Validate(); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	if err := h.usecase.UpdateUserAccess(req.UserId, req.ServiceIds, req.Status); err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			response.NotFound(w, "User access not found", err)
			return
		}
		response.InternalServerError(w, "Failed to update user access", err)
		return
	}

	response.Success(w, "User access updated successfully", nil)
}

func (h *UserAccessHandler) DeleteUserAccess(w http.ResponseWriter, r *http.Request) {
	idStr := request.GetURLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid access ID", err)
		return
	}

	if err := h.usecase.DeleteUserAccess(id); err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			response.NotFound(w, "User access not found", err)
			return
		}
		response.InternalServerError(w, "Failed to delete user access", err)
		return
	}

	response.Success(w, "User access deleted successfully", nil)
}
