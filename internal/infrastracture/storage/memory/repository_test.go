package memory

import (
	"github.com/google/uuid"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/diagnoses"
	"testing"
	"time"
)

func TestRepository_AddDiagnosis(t *testing.T) {
	repo := NewRepository()

	newDiagnosis := diagnoses.Diagnosis{
		ID:           uuid.MustParse("11111111-1111-1111-1111-111111111112"),
		Description:  "test desc",
		PatientID:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		CreatedAt:    time.Now(),
		Prescription: nil,
	}
	err := repo.AddDiagnosis(newDiagnosis)
	if err != nil {
		t.Errorf("AddDiagnosis() Error = %v, but no error expected", err)
	}

	got := repo.diagnoses["11111111-1111-1111-1111-111111111112"]
	if got != newDiagnosis {
		t.Errorf("got=%v, expected=%v", got, newDiagnosis)
	}
}

func TestRepository_GetByID(t *testing.T) {
	repo := NewRepository()

	got, err := repo.GetByID(uuid.MustParse("11111111-1111-1111-1111-111111111111"))
	if got == nil {
		t.Errorf("got <nil>, but a value was expected")
	}

	if err != nil {
		t.Errorf("got error=%v, but no error expected", err)
	}

	nameExpected := "John Doe"
	if got.Name != nameExpected {
		t.Errorf("got=%s, expected John Doe", nameExpected)
	}
}

func TestRepository_GetByName(t *testing.T) {
	repo := NewRepository()
	expected := "John Doe"
	got, err := repo.GetByName(expected)

	if err != nil {
		t.Errorf("got error=%v, but no error expected", err)
	}
	if got == nil {
		t.Errorf("got nil when value was expected")
	}

	if got.Name != expected {
		t.Errorf("got=%s, expected=%s", got.Name, expected)
	}
}
