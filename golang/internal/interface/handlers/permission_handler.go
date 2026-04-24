package handler

import (
	"akastra-access/internal/app/usecases/permission"
	"akastra-access/internal/infrastructure/utils/pagination"
	"akastra-access/internal/interface/http/request"
	"akastra-access/internal/interface/http/response"
	pkgErrors "akastra-access/internal/pkg/errors"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type PermissionHandler struct {
	usecase permission.PermissionUsecase
}

func NewPermissionHandler(u permission.PermissionUsecase) *PermissionHandler {
	return &PermissionHandler{
		usecase: u,
	}
}

func (h *PermissionHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var req request.CreatePermissionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", err)
		return
	}

	if err := req.Validate(); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	res, err := h.usecase.CreatePermission(req.PermissionKey, req.Description)
	if err != nil {
		response.InternalServerError(w, "Failed to create permission", err)
		return
	}

	response.Success(w, "Permission created successfully", res)
}

func (h *PermissionHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	pagination := pagination.GeneratePaginationFromRequest(r)

	permissions, err := h.usecase.GetPermissions(pagination)
	if err != nil {
		response.InternalServerError(w, "Failed to get permissions", err)
		return
	}

	response.Success(w, "Permissions retrieved successfully", permissions)
}

func (h *PermissionHandler) GetPermissionByID(w http.ResponseWriter, r *http.Request) {
	idStr := request.GetURLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid permission ID", err)
		return
	}

	res, err := h.usecase.GetPermissionByID(id)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			response.NotFound(w, "Permission not found", err)
			return
		}
		response.InternalServerError(w, "Failed to get permission", err)
		return
	}

	response.Success(w, "Permission retrieved successfully", res)
}

func (h *PermissionHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	idStr := request.GetURLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid permission ID", err)
		return
	}

	var req request.UpdatePermissionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", err)
		return
	}

	if err := req.Validate(); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	if err := h.usecase.UpdatePermission(id, req.PermissionKey, req.Description); err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			response.NotFound(w, "Permission not found", err)
			return
		}
		response.InternalServerError(w, "Failed to update permission", err)
		return
	}

	response.Success(w, "Permission updated successfully", nil)
}

func (h *PermissionHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	idStr := request.GetURLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid permission ID", err)
		return
	}

	if err := h.usecase.DeletePermission(id); err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			response.NotFound(w, "Permission not found", err)
			return
		}
		response.InternalServerError(w, "Failed to delete permission", err)
		return
	}

	response.Success(w, "Permission deleted successfully", nil)
}
