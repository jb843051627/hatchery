package model

import "fmt"

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s — %s", e.Field, e.Message)
}

var (
	ErrIncubatorNotFound   = fmt.Errorf("incubator not found")
	ErrBatchNotFound       = fmt.Errorf("batch not found")
	ErrEggTrayNotFound     = fmt.Errorf("egg tray not found")
	ErrReadingNotFound     = fmt.Errorf("reading not found")
	ErrScheduleNotFound    = fmt.Errorf("schedule not found")
	ErrHatchRecordNotFound = fmt.Errorf("hatch record not found")
	ErrAlertNotFound       = fmt.Errorf("alert not found")
	ErrMaintenanceNotFound = fmt.Errorf("maintenance task not found")
	ErrInvalidStatus       = fmt.Errorf("invalid status")
	ErrInvalidPhase        = fmt.Errorf("invalid phase")
	ErrDuplicateTrayNumber = fmt.Errorf("duplicate tray number")
	ErrCapacityExceeded    = fmt.Errorf("incubator capacity exceeded")
	ErrBatchNotIncubating  = fmt.Errorf("batch is not in incubating state")
	ErrScheduleConflict    = fmt.Errorf("schedule time conflict")
)
