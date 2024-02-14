package queries

import (
	"errors"
	"github.com/google/uuid"
	"github.com/juanmabaracat/diagnosis-service/internal/app/diagnoses/commands"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/diagnoses"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/patients"
	"reflect"
	"testing"
	"time"
)

func Test_getDiagnoses_Handle(t *testing.T) {
	tests := []struct {
		name        string
		patientRepo patients.Repository
		query       GetDiagnosesQuery
		want        []*diagnoses.Diagnosis
		wantErr     error
	}{
		{
			name: "return error when can't get the patient",
			patientRepo: func() patients.Repository {
				mockRepo := &patients.MockRepository{}
				mockRepo.On("GetByName", "John").Return((*patients.Patient)(nil), errors.New("DB error"))
				return mockRepo
			}(),
			query:   GetDiagnosesQuery{PatientName: "John"},
			want:    nil,
			wantErr: commands.ErrGettingPatient,
		},
		{
			name: "return error when the patient doesn't exists",
			patientRepo: func() patients.Repository {
				mockRepo := &patients.MockRepository{}
				mockRepo.On("GetByName", "John").Return((*patients.Patient)(nil), nil)
				return mockRepo
			}(),
			query:   GetDiagnosesQuery{PatientName: "John"},
			want:    nil,
			wantErr: commands.ErrPatientNotFound,
		},
		{
			name: "return patient diagnoses without error",
			patientRepo: func() patients.Repository {
				mockRepo := &patients.MockRepository{}
				mockRepo.On("GetByName", "John Doe").Return(&patients.Patient{
					ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
					LegalID:     "1234",
					Name:        "Jhon Doe",
					Address:     "test",
					Phone:       "1234",
					Email:       "test@example.com",
					Diagnostics: createFakeDiagnoses(),
				}, nil)
				return mockRepo
			}(),
			query:   GetDiagnosesQuery{PatientName: "John Doe"},
			want:    createFakeDiagnoses(),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &getDiagnoses{patientRepo: tt.patientRepo}
			got, err := g.Handle(tt.query)
			if err != nil && err != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handle() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func createFakeDiagnoses() []*diagnoses.Diagnosis {
	return []*diagnoses.Diagnosis{{
		ID:           uuid.MustParse("11111111-1111-1111-1111-111111111112"),
		Description:  "test diagnosis description",
		PatientID:    uuid.UUID{},
		CreatedAt:    time.Time{},
		Prescription: nil,
	}}
}
