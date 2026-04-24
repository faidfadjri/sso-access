package handler

import (
	"akastra-access/internal/app/usecases/assigned_roles"
	"akastra-access/internal/infrastructure/utils/pagination"
	"akastra-access/internal/interface/http/request"
	"akastra-access/internal/interface/http/response"
	pkgErrors "akastra-access/internal/pkg/errors"
	"encoding/json"
	"errors"
	"net/http"
)

type AssignedRolesHandler struct {
	usecase assigned_roles.AssignedRolesUsecase
}

func NewAssignedRolesHandler(u assigned_roles.AssignedRolesUsecase) *AssignedRolesHandler {
	return &AssignedRolesHandler{
		usecase: u,
	}
}

func (h *AssignedRolesHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	var req request.CreateUserServiceRoleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", err)
		return
	}

	if err := req.Validate(); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	res, err := h.usecase.AssignRole(req.UserIds, req.RoleId)
	if err != nil {
		response.InternalServerError(w, "Failed to assign role to user", err)
		return
	}

	response.Success(w, "Role assigned to user successfully", res)
}

func (h *AssignedRolesHandler) GetList(w http.ResponseWriter, r *http.Request) {
	pagination := pagination.GeneratePaginationFromRequest(r)

	list, err := h.usecase.GetList(pagination)
	if err != nil {
		response.InternalServerError(w, "Failed to get user service roles", err)
		return
	}

	response.Success(w, "User service roles retrieved successfully", list)
}

func (h *AssignedRolesHandler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	// For removing, we can accept Query Params or Body.
	// Body is cleaner for composite keys but GET/DELETE with body is sometimes frowned upon (though standard allows it optionally).
	// Let's support Body for Delete to be consistent with Create for composite keys.

	var req request.DeleteUserServiceRoleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Fallback to query params if body is empty or invalid?
		// Let's try to parse query params if decode fails or if struct is empty.
		// For now strict Body.
		response.BadRequest(w, "Invalid request body", err)
		return
	}

	if err := req.Validate(); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	if err := h.usecase.RemoveRole(req.UserId, req.RoleId); err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			response.NotFound(w, "User service role not found", err)
			return
		}
		response.InternalServerError(w, "Failed to remove role from user", err)
		return
	}

	response.Success(w, "Role removed from user successfully", nil)
}
