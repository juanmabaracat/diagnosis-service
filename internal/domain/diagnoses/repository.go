package diagnoses

type Repository interface {
	AddDiagnosis(diagnosis Diagnosis) error
}
