package commands

import (
	"errors"
	"github.com/google/uuid"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/diagnoses"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/patients"
	"log/slog"
	"time"
)

var (
	ErrPatientNotFound = errors.New("patient not found")
	ErrGettingPatient  = errors.New("error getting patient")
	ErrAddingDiagnosis = errors.New("error adding diagnosis")
	ErrUpdatingPatient = errors.New("error updating patient")
)

type AddPatientDiagnosis struct {
	PatientID    uuid.UUID
	Diagnosis    string
	Prescription *string
}

type AddPatientDiagnosisHandler interface {
	Handle(command AddPatientDiagnosis) error
}

type addPatientDiagnosisHandler struct {
	patientRepo   patients.Repository
	diagnosisRepo diagnoses.Repository
}

func NewAddPatientDiagnosisHandler(patientRepo patients.Repository, diagnosisRepo diagnoses.Repository) AddPatientDiagnosisHandler {
	return &addPatientDiagnosisHandler{
		patientRepo:   patientRepo,
		diagnosisRepo: diagnosisRepo,
	}
}

func (h *addPatientDiagnosisHandler) Handle(command AddPatientDiagnosis) error {
	patient, err := h.patientRepo.GetByID(command.PatientID)
	if err != nil {
		slog.Error(err.Error(), "patientID", command.PatientID)
		return ErrGettingPatient
	}

	if patient == nil {
		slog.Info(ErrPatientNotFound.Error(), "patientID", command.PatientID)
		return ErrPatientNotFound
	}

	newDiagnosis := diagnoses.Diagnosis{
		ID:           uuid.New(),
		Description:  command.Diagnosis,
		PatientID:    patient.ID,
		CreatedAt:    time.Now(),
		Prescription: command.Prescription,
	}

	patient.Diagnostics = append(patient.Diagnostics, &newDiagnosis)

	updateErr := h.patientRepo.Update(*patient)
	if updateErr != nil {
		slog.Error(updateErr.Error(), "patient", *patient)
		return ErrUpdatingPatient
	}

	addErr := h.diagnosisRepo.AddDiagnosis(newDiagnosis)
	if addErr != nil {
		slog.Error(addErr.Error(), "newDiagnosis", newDiagnosis)
		return ErrAddingDiagnosis
	}

	slog.Info("patient diagnosis successfully added", "newDiagnosis", newDiagnosis)
	return nil
}
