package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type SessionRepositoryMock struct {
	mock.Mock
}

func (m *SessionRepositoryMock) SetSession(key string, value interface{}, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func (m *SessionRepositoryMock) GetSession(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *SessionRepositoryMock) DeleteSession(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *SessionRepositoryMock) GetAuthCodePrefix() string {
	args := m.Called()
	return args.String(0)
}

func (m *SessionRepositoryMock) SaveAuthCode(code, sessionID string) error {
	args := m.Called(code, sessionID)
	return args.Error(0)
}

func (m *SessionRepositoryMock) GetAuthCode(code string) (string, error) {
	args := m.Called(code)
	return args.String(0), args.Error(1)
}

func (m *SessionRepositoryMock) DeleteAuthCode(code string) error {
	args := m.Called(code)
	return args.Error(0)
}
