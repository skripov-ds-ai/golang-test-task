package entities

// ResultStatus is a type to express result statuses
type ResultStatus string

const (
	// Success is constant to show success result
	Success ResultStatus = "success"
	// Error is constant to show error result
	Error ResultStatus = "error"
)
