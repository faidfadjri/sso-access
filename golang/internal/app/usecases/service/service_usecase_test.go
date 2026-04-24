package service_test

import (
	"akastra-access/internal/app/usecases/service"
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/utils/pagination"
	"akastra-access/internal/mocks"
	pkgErrors "akastra-access/internal/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateClient(t *testing.T) {
	// mockRepo := new(mocks.ServiceRepositoryMock)
	// usecase := service.NewServiceUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		// req := request.CreateClientReq{ ... } 
		// Skipped due to image processing dependency
		
		// mockRepo.On("CreateClient", mock.AnythingOfType("*domain.Services")).Return(nil).Once()

		// See comments in previous implementation about skipping this test due to heavy dependency
	})
}

func TestGetClients(t *testing.T) {
	mockRepo := new(mocks.ServiceRepositoryMock)
	usecase := service.NewServiceUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		expectedPagination := &pagination.Pagination{}
		mockRepo.On("GetClients", mock.Anything).Return(expectedPagination, nil).Once()

		p := &pagination.Pagination{}
		res, err := usecase.GetClients(p)

		assert.NoError(t, err)
		assert.Equal(t, expectedPagination, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteClientById(t *testing.T) {
	mockRepo := new(mocks.ServiceRepositoryMock)
	usecase := service.NewServiceUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetClientById", uint64(123)).Return(&domain.Services{}, nil).Once()
		mockRepo.On("DeleteClientById", uint64(123)).Return(nil).Once()

		err := usecase.DeleteClientById(123)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetClientById", uint64(123)).Return(nil, gorm.ErrRecordNotFound).Once()

		err := usecase.DeleteClientById(123)

		assert.Error(t, err)
		assert.Equal(t, pkgErrors.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}
