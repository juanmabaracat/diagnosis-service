package commands

import (
	"errors"
	"github.com/google/uuid"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/diagnoses"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/patients"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_addPatientDiagnosisHandler_Handle(t *testing.T) {
	patientID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	command := AddPatientDiagnosis{
		PatientID:    patientID,
		Diagnosis:    "test diagnosis",
		Prescription: nil,
	}

	tests := []struct {
		name          string
		patientRepo   patients.Repository
		diagnosisRepo diagnoses.Repository
		command       AddPatientDiagnosis
		wantErr       error
	}{
		{
			name: "return error when fails getting patient",
			patientRepo: func() patients.Repository {
				mockRepo := &patients.MockRepository{}
				mockRepo.On("GetByID", patientID).Return((*patients.Patient)(nil), errors.New("cannot get user"))
				return mockRepo
			}(),
			diagnosisRepo: &diagnoses.MockRepository{},
			command:       command,
			wantErr:       ErrGettingPatient,
		},
		{
			name: "return error when there is no patient for that ID",
			patientRepo: func() patients.Repository {
				mockRepo := &patients.MockRepository{}
				mockRepo.On("GetByID", patientID).Return((*patients.Patient)(nil), nil)
				return mockRepo
			}(),
			diagnosisRepo: &diagnoses.MockRepository{},
			command:       command,
			wantErr:       ErrPatientNotFound,
		},
		{
			name: "return error when the patient cant be updated",
			patientRepo: func() patients.Repository {
				mockRepo := &patients.MockRepository{}
				patient := &patients.Patient{}
				mockRepo.On("GetByID", patientID).Return(patient, nil)
				mockRepo.On("Update", mock.Anything).Return(errors.New("update error"))
				return mockRepo
			}(),
			diagnosisRepo: &diagnoses.MockRepository{},
			command:       command,
			wantErr:       ErrUpdatingPatient,
		},
		{
			name: "return error when the patient cant be updated",
			patientRepo: func() patients.Repository {
				mockRepo := &patients.MockRepository{}
				patient := &patients.Patient{}
				mockRepo.On("GetByID", patientID).Return(patient, nil)
				mockRepo.On("Update", mock.Anything).Return(nil)
				return mockRepo
			}(),
			diagnosisRepo: func() diagnoses.Repository {
				mockRepo := &diagnoses.MockRepository{}
				mockRepo.On("AddDiagnosis", mock.Anything).Return(errors.New("add error"))
				return mockRepo
			}(),
			command: command,
			wantErr: ErrAddingDiagnosis,
		},
		{
			name: "add patient diagnosis without error",
			patientRepo: func() patients.Repository {
				mockRepo := &patients.MockRepository{}
				patient := &patients.Patient{}
				mockRepo.On("GetByID", patientID).Return(patient, nil)
				mockRepo.On("Update", mock.Anything).Return(nil)
				return mockRepo
			}(),
			diagnosisRepo: func() diagnoses.Repository {
				mockRepo := &diagnoses.MockRepository{}
				mockRepo.On("AddDiagnosis", mock.Anything).Return(nil)
				return mockRepo
			}(),
			command: command,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &addPatientDiagnosisHandler{
				patientRepo:   tt.patientRepo,
				diagnosisRepo: tt.diagnosisRepo,
			}
			if err := h.Handle(tt.command); !errors.Is(err, tt.wantErr) {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
