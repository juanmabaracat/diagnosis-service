package commands

import "github.com/stretchr/testify/mock"

type MockAddPatientDiagnosis struct {
	mock.Mock
}

func (m *MockAddPatientDiagnosis) Handle(command AddPatientDiagnosis) error {
	args := m.Called(command)
	return args.Error(0)
}
