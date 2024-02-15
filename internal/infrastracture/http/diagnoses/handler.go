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

// AddDiagnosis godoc
//
//	@Summary		Add patient diagnosis
//	@Description	Add patient diagnosis
//	@Tags			diagnosis
//	@Accept			json
//	@Produce		json
//	@Param			patientID			path		string		true	"patient ID"
//	@Param			diagnosis body		AddDiagnosisRequest		true	"add diagnosis"
//	@Success		201	{string}		status created
//	@Failure		400	{object}		HTTPError
//	@Failure		404	{object}		HTTPError
//	@Failure		500	{object}		HTTPError
//	@Router			/patient/{patientID}/diagnoses [post]
func (h *Handler) AddDiagnosis(writer http.ResponseWriter, request *http.Request) {
	addDiagnosisRequest := AddDiagnosisRequest{}
	patientIDParam := chi.URLParam(request, PatientIDURLParam)
	patientID, parseErr := uuid.Parse(patientIDParam)
	if parseErr != nil {
		writeError(writer, http.StatusBadRequest, errInvalidID)
		return
	}

	decodeErr := json.NewDecoder(request.Body).Decode(&addDiagnosisRequest)
	if decodeErr != nil {
		writeError(writer, http.StatusBadRequest, decodeErr)
		return
	}

	addDiagnosisRequest.Diagnosis = strings.TrimSpace(addDiagnosisRequest.Diagnosis)
	if addDiagnosisRequest.Diagnosis == "" {
		writeError(writer, http.StatusBadRequest, errInvalidDiagnosis)
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
			writeError(writer, http.StatusNotFound, errPatientNotFound)
			return
		}
		writeError(writer, http.StatusInternalServerError, errProcessingRequest)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	return
}

type GetDiagnosesResponse struct {
	PatientName string                 `json:"patient_name"`
	Diagnoses   []*diagnoses.Diagnosis `json:"patient_diagnoses"`
}

// GetDiagnoses godoc
//
//	@Summary		Get patient diagnoses
//	@Description	Get patient diagnoses
//	@Tags			diagnosis
//	@Accept			json
//	@Produce		json
//	@Param			patientName				query					string	true	"diagnoses search by patient name"
//	@Success		200	{object}			GetDiagnosesResponse
//	@Failure		400	{object}			HTTPError
//	@Failure		404	{object}			HTTPError
//	@Failure		500	{object}			HTTPError
//	@Router			/patient/diagnoses 		[get]
func (h *Handler) GetDiagnoses(writer http.ResponseWriter, request *http.Request) {
	patientName := request.URL.Query().Get(PatientNameQueryParam)
	patientName = strings.TrimSpace(patientName)
	if patientName == "" {
		writeError(writer, http.StatusBadRequest, errInvalidPatientName)
		return
	}

	pDiagnoses, err := h.diagnosesServices.Queries.GetDiagnoses.Handle(queries.GetDiagnosesQuery{PatientName: patientName})
	if err != nil {
		if errors.Is(err, commands.ErrPatientNotFound) {
			slog.Info("patient not found", "patientName", patientName)
			writeError(writer, http.StatusNotFound, errPatientNotFound)
			return
		}

		slog.Error("error getting diagnoses", "err", err)
		writeError(writer, http.StatusInternalServerError, errProcessingRequest)
		return
	}

	encodeErr := json.NewEncoder(writer).Encode(GetDiagnosesResponse{
		PatientName: patientName,
		Diagnoses:   pDiagnoses,
	})
	if encodeErr != nil {
		slog.Error("error encoding get diagnoses response", "encodeErr", encodeErr)
		writeError(writer, http.StatusInternalServerError, errProcessingRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
	return
}

func writeError(writer http.ResponseWriter, code int, err error) {
	writer.WriteHeader(code)
	httpErr := HTTPError{
		Code:    code,
		Message: err.Error(),
	}
	errEncode := json.NewEncoder(writer).Encode(httpErr)
	if errEncode != nil {
		slog.Error("error encoding http error", "err", errEncode)
		return
	}
}

// HTTP HTTPError
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
