package handler

import (
	"akastra-access/internal/app/usecases/role"
	"akastra-access/internal/infrastructure/utils/pagination"
	"akastra-access/internal/interface/http/request"
	"akastra-access/internal/interface/http/response"
	pkgErrors "akastra-access/internal/pkg/errors"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type RoleHandler struct {
	usecase role.RoleUsecase
}

func NewRoleHandler(u role.RoleUsecase) *RoleHandler {
	return &RoleHandler{
		usecase: u,
	}
}

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req request.CreateRoleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", err)
		return
	}

	if err := req.Validate(); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	res, err := h.usecase.CreateRole(req.RoleName, req.ServiceId)
	if err != nil {
		response.InternalServerError(w, "Failed to create role", err)
		return
	}

	response.Success(w, "Role created successfully", res)
}

func (h *RoleHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	pagination := pagination.GeneratePaginationFromRequest(r)

	roles, err := h.usecase.GetRoles(pagination)
	if err != nil {
		response.InternalServerError(w, "Failed to get roles", err)
		return
	}

	response.Success(w, "Roles retrieved successfully", roles)
}

func (h *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	idStr := request.GetURLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid role ID", err)
		return
	}

	role, err := h.usecase.GetRoleByID(id)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			response.NotFound(w, "Role not found", err)
			return
		}
		response.InternalServerError(w, "Failed to get role", err)
		return
	}

	response.Success(w, "Role retrieved successfully", role)
}

func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	idStr := request.GetURLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid role ID", err)
		return
	}

	var req request.UpdateRoleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", err)
		return
	}

	if err := req.Validate(); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	if err := h.usecase.UpdateRole(id, req.RoleName, req.ServiceId); err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			response.NotFound(w, "Role not found", err)
			return
		}
		response.InternalServerError(w, "Failed to update role", err)
		return
	}

	response.Success(w, "Role updated successfully", nil)
}

func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	idStr := request.GetURLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid role ID", err)
		return
	}

	if err := h.usecase.DeleteRole(id); err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			response.NotFound(w, "Role not found", err)
			return
		}
		response.InternalServerError(w, "Failed to delete role", err)
		return
	}

	response.Success(w, "Role deleted successfully", nil)
}
