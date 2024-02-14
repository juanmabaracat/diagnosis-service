package patients

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetByName(name string) (*Patient, error) {
	args := m.Called(name)
	return args.Get(0).(*Patient), args.Error(1)
}

func (m *MockRepository) GetByID(ID uuid.UUID) (*Patient, error) {
	args := m.Called(ID)
	return args.Get(0).(*Patient), args.Error(1)
}

func (m *MockRepository) Update(patient Patient) error {
	args := m.Called(patient)
	return args.Error(0)
}
