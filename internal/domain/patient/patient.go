package patient

import (
	"github.com/google/uuid"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/diagnosis"
)

type Patient struct {
	ID          uuid.UUID
	LegalID     string
	Name        string
	Address     string
	Phone       string
	Email       string
	diagnostics []*diagnosis.Diagnosis
}
