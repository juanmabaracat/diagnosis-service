package diagnoses

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juanmabaracat/diagnosis-service/internal/app"
	"github.com/juanmabaracat/diagnosis-service/internal/app/diagnoses/commands"
	"github.com/juanmabaracat/diagnosis-service/internal/app/diagnoses/queries"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/diagnoses"
	"log/slog"
	"net/http"
	"strings"
)

var (
	errInvalidID          = errors.New("invalid ID")
	errInvalidDiagnosis   = errors.New("diagnosis cannot be empty")
	errPatientNotFound    = errors.New("there no patient for the ID supplied")
	errProcessingRequest  = errors.New("error processing the request")
	errInvalidPatientName = errors.New("invalid patient name")
)

const (
	PatientIDURLParam     = "patientID"
	PatientNameQueryParam = "patientName"
)

type Handler struct {
	diagnosesServices app.DiagnosisServices
}

func NewHandler(diagnosisServices app.DiagnosisServices) *Handler {
	return &Handler{
		diagnosesServices: diagnosisServices,
	}
}

type AddDiagnosisRequest struct {
	Diagnosis    string  `json:"diagnosis"`
	Prescription *string `json:"prescription"`
}

func (h *Handler) AddDiagnosis(writer http.ResponseWriter, request *http.Request) {
	addDiagnosisRequest := AddDiagnosisRequest{}
	patientIDParam := chi.URLParam(request, PatientIDURLParam)
	patientID, parseErr := uuid.Parse(patientIDParam)
	if parseErr != nil {
		http.Error(writer, errInvalidID.Error(), http.StatusBadRequest)
		return
	}

	decodeErr := json.NewDecoder(request.Body).Decode(&addDiagnosisRequest)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), http.StatusBadRequest)
		return
	}

	addDiagnosisRequest.Diagnosis = strings.TrimSpace(addDiagnosisRequest.Diagnosis)
	if addDiagnosisRequest.Diagnosis == "" {
		http.Error(writer, errInvalidDiagnosis.Error(), http.StatusBadRequest)
		return
	}

	err := h.diagnosesServices.Commands.AddPatientDiagnosisHandler.Handle(commands.AddPatientDiagnosis{
		PatientID:    patientID,
		Diagnosis:    addDiagnosisRequest.Diagnosis,
		Prescription: addDiagnosisRequest.Prescription,
	})

	if err != nil {
		slog.Error("error handling request for adding diagnosis", "error", err)
		if errors.Is(err, commands.ErrPatientNotFound) {
			http.Error(writer, errPatientNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(writer, errProcessingRequest.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	return
}

type GetDiagnosesResponse struct {
	PatientName string                 `json:"patient_name"`
	Diagnoses   []*diagnoses.Diagnosis `json:"patient_diagnoses"`
}

func (h *Handler) GetDiagnoses(writer http.ResponseWriter, request *http.Request) {
	patientName := request.URL.Query().Get(PatientNameQueryParam)
	patientName = strings.TrimSpace(patientName)
	if patientName == "" {
		http.Error(writer, errInvalidPatientName.Error(), http.StatusBadRequest)
		return
	}

	pDiagnoses, err := h.diagnosesServices.Queries.GetDiagnoses.Handle(queries.GetDiagnosesQuery{PatientName: patientName})
	if err != nil {
		if errors.Is(err, commands.ErrPatientNotFound) {
			slog.Info("patient not found", "patientName", patientName)
			http.Error(writer, errPatientNotFound.Error(), http.StatusNotFound)
			return
		}
		slog.Error("error getting diagnoses", "err", err)
		http.Error(writer, errProcessingRequest.Error(), http.StatusInternalServerError)
		return
	}

	encodeErr := json.NewEncoder(writer).Encode(GetDiagnosesResponse{
		PatientName: patientName,
		Diagnoses:   pDiagnoses,
	})
	if encodeErr != nil {
		slog.Error("error encoding get diagnoses response", "encodeErr", encodeErr)
		http.Error(writer, errProcessingRequest.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	return
}
