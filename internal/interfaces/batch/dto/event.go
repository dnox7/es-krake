package dto

import "time"

type BatchEvent struct {
	ID         string                 `header:"X-Event-Id"`                                           // Unique identifier for the event
	Timestamp  time.Time              `header:"X-Event-Time" time_format:"2006-01-02T15:04:05Z07:00"` // Event creation or received time
	Type       string                 `header:"X-Event-Type"`                                         // Type or category of the event
	BatchLogID int                    `header:"X-Batch-Log-Id"`                                       // Internal log or DB reference
	Data       map[string]interface{} // Event payload from request body
}
