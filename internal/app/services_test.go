package app

import (
	"github.com/juanmabaracat/diagnosis-service/internal/app/diagnoses/commands"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/diagnoses"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/patients"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServices(t *testing.T) {
	patientRepo := &patients.MockRepository{}
	diagnosisRepo := &diagnoses.MockRepository{}
	expected := Services{
		DiagnosisServices: DiagnosisServices{
			Commands: Commands{
				AddPatientDiagnosisHandler: commands.NewAddPatientDiagnosisHandler(patientRepo, diagnosisRepo),
			},
		},
	}

	got := NewServices(patientRepo, diagnosisRepo)

	assert.Equal(t, got, expected)
}
