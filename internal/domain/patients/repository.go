package patients

import "github.com/google/uuid"

type Repository interface {
	GetByName(name string) (*Patient, error)
	GetByID(ID uuid.UUID) (*Patient, error)
	Update(patient Patient) error
}
