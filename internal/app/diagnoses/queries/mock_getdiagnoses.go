package queries

import (
	"github.com/juanmabaracat/diagnosis-service/internal/domain/diagnoses"
	"github.com/stretchr/testify/mock"
)

type MockGetDiagnoses struct {
	mock.Mock
}

func (m *MockGetDiagnoses) Handle(query GetDiagnosesQuery) ([]*diagnoses.Diagnosis, error) {
	args := m.Called(query)
	return args.Get(0).([]*diagnoses.Diagnosis), args.Error(1)
}
