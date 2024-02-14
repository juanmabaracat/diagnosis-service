package diagnoses

import "github.com/stretchr/testify/mock"

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) AddDiagnosis(diagnosis Diagnosis) error {
	args := m.Called(diagnosis)
	return args.Error(0)
}
