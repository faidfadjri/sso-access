package handler

import (
	"akastra-access/internal/app/usecases"
	"akastra-access/internal/infrastructure/utils/pagination"
	"akastra-access/internal/interface/http/request"
	"akastra-access/internal/interface/http/response"
	pkgErrors "akastra-access/internal/pkg/errors"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// ServiceHandler handles HTTP requests
type ServiceHandler struct {
	usecase usecases.ServiceUsecase
}

func NewServiceHandler(u usecases.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{
		usecase: u,
	}
}

func (h *ServiceHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	// parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		response.BadRequest(w, "Invalid multipart form", err)
		return
	}

	// get file
	file, header, err := r.FormFile("logo")
	if err != nil {
		response.BadRequest(w, "Logo file is required", err)
		return
	}
	defer file.Close()

	// get request body
	req := request.CreateClientReq{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		RedirectURL: r.FormValue("redirect_url"),
	}

	// validate request body
	if err := req.Validate(); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	// create client
	service, err := h.usecase.CreateClient(req, file, header)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrInvalidFormat) {
			response.BadRequest(w, err.Error(), nil)
			return
		}
		log.Println(err.Error())
		response.InternalServerError(w, "Failed to create client", nil)
		return
	}

	// success response
	response.Success(w, "Client created successfully", service)
}

func (h *ServiceHandler) DeleteClientById(w http.ResponseWriter, r *http.Request) {
	// get service id
	serviceId := request.GetURLParam(r, "id")
	serviceIdUint64, err := strconv.ParseUint(serviceId, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid service id", err)
		return
	}

	// delete client
	if err := h.usecase.DeleteClientById(serviceIdUint64); err != nil {
		log.Println(err.Error())

		if errors.Is(err, pkgErrors.ErrNotFound) {
			response.NotFound(w, "Client not found", err)
			return
		}

		response.InternalServerError(w, "Failed to delete client", err)
		return
	}

	// success response
	response.Success(w, "Client deleted successfully", nil)
}

func (h *ServiceHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	// service_id
	serviceId := request.GetURLParam(r, "id")
	serviceIdUint64, err := strconv.ParseUint(serviceId, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid service id", err)
		return
	}

	// parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		response.BadRequest(w, "Invalid multipart form", err)
		return
	}

	// get file
	file, header, _ := r.FormFile("logo")
	if file != nil {
		defer file.Close()
	}


	isActiveStr := r.FormValue("is_active")

	var isActivePtr *bool
	if isActiveStr != "" {
		parsed, err := strconv.ParseBool(isActiveStr)
		if err != nil {
			response.BadRequest(w, "Invalid is_active value", err)
			return
		}
		isActivePtr = &parsed
	}

	// get request body
	req := request.UpdateClientReq{
		ID:          serviceIdUint64,
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		RedirectURL: r.FormValue("redirect_url"),
		IsActive:    isActivePtr,
	}

	// validate request body
	if err := req.Validate(); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	// update client
	if err := h.usecase.UpdateClient(req, file, header); err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			response.NotFound(w, "Client not found", err)
			return
		}
		log.Println(err.Error())
		response.InternalServerError(w, "Failed to update client", nil)
		return
	}

	// success response
	response.Success(w, "Client updated successfully", nil)
}

func (h *ServiceHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	// generate pagination from request
	pagination := pagination.GeneratePaginationFromRequest(r)

	// get clients
	clients, err := h.usecase.GetClients(pagination)
	if err != nil {
		response.InternalServerError(w, "Failed to get clients", err)
		return
	}

	// success response
	response.Success(w, "Clients retrieved successfully", clients)
}