package patients

import (
	"github.com/google/uuid"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/diagnoses"
)

type Patient struct {
	ID          uuid.UUID
	LegalID     string
	Name        string
	Address     string
	Phone       string
	Email       string
	Diagnostics []*diagnoses.Diagnosis
}
