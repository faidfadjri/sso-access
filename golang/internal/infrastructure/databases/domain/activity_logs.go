package domain

import (
	"time"
)

type ActivityLogs struct {
	LogId        uint64    `gorm:"primaryKey;column:log_id" json:"log_id"`
	UserId       uint64    `gorm:"column:user_id" json:"user_id"`
	Action       string    `gorm:"column:action;size:50" json:"action"`
	ChangeFields string    `gorm:"column:change_fields;type:json" json:"change_fields"` // or []byte or custom JSON type
	Before       string    `gorm:"column:before;type:json" json:"before"`
	After        string    `gorm:"column:after;type:json" json:"after"`
	Desc         string    `gorm:"column:desc;type:text" json:"desc"`
	AffectedTableName string    `gorm:"column:table_name;size:100" json:"table_name"`
	RecordId          int64     `gorm:"column:record_id" json:"record_id"`
	CanUndo           bool      `gorm:"column:can_undo;default:true" json:"can_undo"`
	UndoneBy          *uint64   `gorm:"column:undone_by" json:"undone_by"`
	IpAddress         *string   `gorm:"column:ip_address;size:100" json:"ip_address"`
	UndoneAt          *time.Time `gorm:"column:undone_at" json:"undone_at"`
	BaseModel
}

func (ActivityLogs) TableName() string {
	return "activity_logs"
}
