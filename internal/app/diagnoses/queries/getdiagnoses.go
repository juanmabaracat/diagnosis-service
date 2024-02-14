package queries

import (
	"github.com/juanmabaracat/diagnosis-service/internal/app/diagnoses/commands"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/diagnoses"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/patients"
	"log/slog"
)

type GetDiagnosesQuery struct {
	PatientName string
}

type GetDiagnosesHandler interface {
	Handle(query GetDiagnosesQuery) ([]*diagnoses.Diagnosis, error)
}

type getDiagnoses struct {
	patientRepo patients.Repository
}

func NewGetDiagnosesHandler(patientRepo patients.Repository) GetDiagnosesHandler {
	return &getDiagnoses{patientRepo: patientRepo}
}

func (g *getDiagnoses) Handle(query GetDiagnosesQuery) ([]*diagnoses.Diagnosis, error) {
	patient, err := g.patientRepo.GetByName(query.PatientName)
	if err != nil {
		slog.Error("error getting patient", "err", err, "query", query)
		return nil, commands.ErrGettingPatient
	}

	if patient == nil {
		return nil, commands.ErrPatientNotFound
	}

	return patient.Diagnostics, nil
}
