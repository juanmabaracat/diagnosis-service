package diagnosis

import (
	"github.com/google/uuid"
	"time"
)

type Diagnosis struct {
	ID           uuid.UUID
	Description  string
	PatientID    uuid.UUID
	CreatedAt    time.Time
	Prescription *Prescription
}

type Prescription struct {
	Description string
	Medicine    string
	Doses       string
}
