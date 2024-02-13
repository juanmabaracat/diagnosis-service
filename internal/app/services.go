package app

import (
	"github.com/juanmabaracat/diagnosis-service/internal/app/diagnoses/commands"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/diagnoses"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/patients"
)

type Commands struct {
	AddPatientDiagnosisHandler commands.AddPatientDiagnosisHandler
}

type DiagnosisServices struct {
	Commands Commands
}

// Services contains all services exposed of the application layer
type Services struct {
	DiagnosisServices DiagnosisServices
}

func NewServices(patientRepo patients.Repository, diagnosisRepo diagnoses.Repository) Services {
	return Services{
		DiagnosisServices: DiagnosisServices{
			Commands: Commands{
				AddPatientDiagnosisHandler: commands.NewAddPatientDiagnosisHandler(patientRepo, diagnosisRepo),
			},
		},
	}
}
