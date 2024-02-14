package diagnoses

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juanmabaracat/diagnosis-service/internal/app"
	"github.com/juanmabaracat/diagnosis-service/internal/app/diagnoses/commands"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_AddDiagnosis(t *testing.T) {
	tests := []struct {
		name       string
		handler    commands.AddPatientDiagnosisHandler
		body       interface{}
		PatientID  string
		wantStatus int
		wantMsg    string
	}{
		{
			name:    "return bad request on invalid patient id",
			handler: nil,
			body: AddDiagnosisRequest{
				Diagnosis:    "test diagnosis",
				Prescription: nil,
			},
			PatientID:  "",
			wantStatus: 400,
			wantMsg:    errInvalidID.Error(),
		},
		{
			name:    "return bad request on invalid diagnosis",
			handler: nil,
			body: AddDiagnosisRequest{
				Diagnosis:    "   \n    ",
				Prescription: nil,
			},
			PatientID:  "11111111-1111-1111-1111-111111111111",
			wantStatus: 400,
			wantMsg:    errInvalidDiagnosis.Error(),
		},
		{
			name: "return not found when the patient ID doesn't exists",
			handler: func() commands.AddPatientDiagnosisHandler {
				mock := &commands.MockAddPatientDiagnosis{}
				mock.On("Handle", commands.AddPatientDiagnosis{
					PatientID:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
					Diagnosis:    "test diagnosis",
					Prescription: nil,
				}).Return(commands.ErrPatientNotFound)
				return mock
			}(),
			body: AddDiagnosisRequest{
				Diagnosis:    "test diagnosis",
				Prescription: nil,
			},
			PatientID:  "11111111-1111-1111-1111-111111111111",
			wantStatus: 404,
			wantMsg:    errPatientNotFound.Error(),
		},
		{
			name: "return server error when there is a error adding the diagnosis",
			handler: func() commands.AddPatientDiagnosisHandler {
				mock := &commands.MockAddPatientDiagnosis{}
				mock.On("Handle", commands.AddPatientDiagnosis{
					PatientID:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
					Diagnosis:    "test diagnosis",
					Prescription: nil,
				}).Return(commands.ErrAddingDiagnosis)
				return mock
			}(),
			body: AddDiagnosisRequest{
				Diagnosis:    "test diagnosis",
				Prescription: nil,
			},
			PatientID:  "11111111-1111-1111-1111-111111111111",
			wantStatus: 500,
			wantMsg:    errProcessingRequest.Error(),
		},
		{
			name: "create the diagnosis without error",
			handler: func() commands.AddPatientDiagnosisHandler {
				mock := &commands.MockAddPatientDiagnosis{}
				mock.On("Handle", commands.AddPatientDiagnosis{
					PatientID:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
					Diagnosis:    "test diagnosis",
					Prescription: nil,
				}).Return(nil)
				return mock
			}(),
			body: AddDiagnosisRequest{
				Diagnosis:    "test diagnosis",
				Prescription: nil,
			},
			PatientID:  "11111111-1111-1111-1111-111111111111",
			wantStatus: 201,
			wantMsg:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(app.DiagnosisServices{Commands: app.Commands{AddPatientDiagnosisHandler: tt.handler}})
			buf := new(bytes.Buffer)
			_ = json.NewEncoder(buf).Encode(tt.body)
			r, _ := http.NewRequest("POST", "/patients/"+tt.PatientID+"/diagnoses", buf)
			rCtx := chi.NewRouteContext()
			rCtx.URLParams.Add(PatientIDURLParam, tt.PatientID)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rCtx))
			response := httptest.NewRecorder()
			h.AddDiagnosis(response, r)

			assert.Equal(t, tt.wantStatus, response.Code)
			assert.Equal(t, tt.wantMsg, strings.TrimRight(response.Body.String(), "\n"))
		})
	}

}
