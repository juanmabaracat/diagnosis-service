package memory

import (
	"github.com/google/uuid"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/diagnoses"
	"github.com/juanmabaracat/diagnosis-service/internal/domain/patients"
	"sync"
)

func NewRepository() Repository {
	repo := Repository{}

	repo.patients = createFakePatients()
	repo.diagnoses = make(map[string]diagnoses.Diagnosis)
	repo.mutex = &sync.RWMutex{}

	return repo
}

type Repository struct {
	patients  map[string]patients.Patient
	diagnoses map[string]diagnoses.Diagnosis
	mutex     *sync.RWMutex
}

func (r *Repository) GetByName(name string) (*patients.Patient, error) {
	for _, p := range r.patients {
		if p.Name == name {
			return &p, nil
		}
	}

	return nil, nil
}

func (r *Repository) GetByID(ID uuid.UUID) (*patients.Patient, error) {
	r.mutex.RLock()
	patient, ok := r.patients[ID.String()]
	r.mutex.RUnlock()
	if !ok {
		return nil, nil
	}

	return &patient, nil
}

func (r *Repository) Update(patient patients.Patient) error {
	r.mutex.Lock()
	r.patients[patient.ID.String()] = patient
	r.mutex.Unlock()
	return nil
}

func (r *Repository) AddDiagnosis(diagnosis diagnoses.Diagnosis) error {
	r.mutex.Lock()
	r.diagnoses[diagnosis.ID.String()] = diagnosis
	r.mutex.Unlock()
	return nil
}

func createFakePatients() map[string]patients.Patient {
	patientID := "11111111-1111-1111-1111-111111111111"

	return map[string]patients.Patient{patientID: {
		ID:          uuid.MustParse(patientID),
		LegalID:     "ABC1234",
		Name:        "John Doe",
		Address:     "Wall Street 123",
		Phone:       "123456789",
		Email:       "john.doe@example.com",
		Diagnostics: []*diagnoses.Diagnosis{},
	}}
}
