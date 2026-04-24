package handler

import (
	"akastra-access/internal/app/usecases/role_permission"
	"akastra-access/internal/infrastructure/utils/pagination"
	"akastra-access/internal/interface/http/request"
	"akastra-access/internal/interface/http/response"
	pkgErrors "akastra-access/internal/pkg/errors"
	"encoding/json"
	"errors"
	"net/http"
)

type RolePermissionHandler struct {
	usecase role_permission.RolePermissionUsecase
}

func NewRolePermissionHandler(u role_permission.RolePermissionUsecase) *RolePermissionHandler {
	return &RolePermissionHandler{
		usecase: u,
	}
}

func (h *RolePermissionHandler) AssignPermission(w http.ResponseWriter, r *http.Request) {
	var req request.CreateRolePermissionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", err)
		return
	}

	if err := req.Validate(); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	res, err := h.usecase.AssignPermission(req.RoleId, req.PermissionId)
	if err != nil {
		response.InternalServerError(w, "Failed to assign permission to role", err)
		return
	}

	response.Success(w, "Permission assigned to role successfully", res)
}

func (h *RolePermissionHandler) GetList(w http.ResponseWriter, r *http.Request) {
	pagination := pagination.GeneratePaginationFromRequest(r)

	list, err := h.usecase.GetList(pagination)
	if err != nil {
		response.InternalServerError(w, "Failed to get role permissions", err)
		return
	}

	response.Success(w, "Role permissions retrieved successfully", list)
}

func (h *RolePermissionHandler) RemovePermission(w http.ResponseWriter, r *http.Request) {
	var req request.DeleteRolePermissionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", err)
		return
	}

	if err := req.Validate(); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	if err := h.usecase.RemovePermission(req.RoleId, req.PermissionId); err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			response.NotFound(w, "Role permission not found", err)
			return
		}
		response.InternalServerError(w, "Failed to remove permission from role", err)
		return
	}

	response.Success(w, "Permission removed from role successfully", nil)
}
