package entity

import (
	"time"

	"github.com/dpe27/es-krake/internal/domain/shared/model"
)

const BatchLogTableName = "batch_logs"

type BatchLog struct {
	ID             int        `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	BatchLogTypeID int        `gorm:"column:batch_log_type_id;type:bigint;not null"           json:"batch_log_type_id"`
	EventID        int        `gorm:"column:event_id;type:bigint;not null"                    json:"event_id"`
	Arguments      string     `gorm:"column:arguments;type:json"                              json:"arguments"`
	StartedAt      time.Time  `gorm:"column:started_at;type:timestamp;not null"               json:"started_at"`
	EndedAt        *time.Time `gorm:"column:ended_at;type:timestamp"                          json:"ended_at"`
	Success        *bool      `gorm:"column:success;type:boolean"                             json:"success"`
	ErrorMessage   *string    `gorm:"column:error_message;type:text"                          json:"error_message"`
	model.BaseModel
}

func (BatchLog) TableName() string {
	return BatchLogTableName
}
