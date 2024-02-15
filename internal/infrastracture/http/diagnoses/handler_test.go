package diagnoses

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juanmabaracat/diagnosis-service/internal/app"
	"github.com/juanmabaracat/diagnosis-service/internal/app/diagnoses/commands"
	"github.com/juanmabaracat/diagnosis-service/internal/app/diagnoses/queries"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/diagnoses"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_AddDiagnosis(t *testing.T) {
	tests := []struct {
		name       string
		handler    commands.AddPatientDiagnosisHandler
		body       interface{}
		PatientID  string
		wantStatus int
		wantErr    *HTTPError
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
			wantErr: &HTTPError{
				Code:    400,
				Message: errInvalidID.Error(),
			},
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
			wantErr: &HTTPError{
				Code:    400,
				Message: errInvalidDiagnosis.Error(),
			},
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
			wantErr: &HTTPError{
				Code:    404,
				Message: errPatientNotFound.Error(),
			},
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
			wantErr: &HTTPError{
				Code:    500,
				Message: errProcessingRequest.Error(),
			},
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
			wantErr:    nil,
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
			if tt.wantErr != nil {
				respErr := HTTPError{}
				err := json.NewDecoder(response.Body).Decode(&respErr)
				assert.Nil(t, err)
				assert.Equal(t, *tt.wantErr, respErr)
			}
		})
	}
}

func TestHandler_GetDiagnoses(t *testing.T) {
	tests := []struct {
		name       string
		queryParam string
		handler    queries.GetDiagnosesHandler
		wantStatus int
	}{
		{
			name:       "return bad request when patient name is invalid",
			queryParam: "patientName=",
			handler: func() queries.GetDiagnosesHandler {
				mock := &queries.MockGetDiagnoses{}
				mock.On("Handle", queries.GetDiagnosesQuery{PatientName: ""}).
					Return(([]*diagnoses.Diagnosis)(nil), commands.ErrGettingPatient)
				return mock
			}(),
			wantStatus: 400,
		},
		{
			name:       "return not found when the patient doesn't exists",
			queryParam: "patientName=John Doe",
			handler: func() queries.GetDiagnosesHandler {
				mock := &queries.MockGetDiagnoses{}
				mock.On("Handle", queries.GetDiagnosesQuery{PatientName: "John Doe"}).
					Return(([]*diagnoses.Diagnosis)(nil), commands.ErrPatientNotFound)
				return mock
			}(),
			wantStatus: 404,
		},
		{
			name:       "return server error when there is an error getting the diagnoses",
			queryParam: "patientName=John Doe",
			handler: func() queries.GetDiagnosesHandler {
				mock := &queries.MockGetDiagnoses{}
				mock.On("Handle", queries.GetDiagnosesQuery{PatientName: "John Doe"}).
					Return(([]*diagnoses.Diagnosis)(nil), commands.ErrGettingPatient)
				return mock
			}(),
			wantStatus: 500,
		},
		{
			name:       "return the diagnoses without error",
			queryParam: "patientName=John Doe",
			handler: func() queries.GetDiagnosesHandler {
				mock := &queries.MockGetDiagnoses{}
				mock.On("Handle", queries.GetDiagnosesQuery{PatientName: "John Doe"}).
					Return([]*diagnoses.Diagnosis{}, nil)
				return mock
			}(),
			wantStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				diagnosesServices: app.DiagnosisServices{
					Queries: app.Queries{GetDiagnoses: tt.handler}},
			}

			req, _ := http.NewRequest("GET", "/patients/diagnoses?"+tt.queryParam, nil)
			resp := httptest.NewRecorder()
			h.GetDiagnoses(resp, req)
			assert.Equal(t, tt.wantStatus, resp.Code)
		})
	}
}
