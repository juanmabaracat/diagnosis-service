package diagnoses

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juanmabaracat/diagnosis-service/internal/app"
	"github.com/juanmabaracat/diagnosis-service/internal/app/diagnoses/commands"
	"log/slog"
	"net/http"
	"strings"
)

var (
	errInvalidID         = errors.New("invalid ID")
	errInvalidDiagnosis  = errors.New("diagnosis cannot be empty")
	errPatientNotFound   = errors.New("there no patient for the ID supplied")
	errProcessingRequest = errors.New("error processing the request")
)

const PatientIDURLParam = "patientID"

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
