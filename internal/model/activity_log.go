package model

import "time"

type ActivityLog struct {
	ID        int `gorm:"primaryKey"`
	Request   string
	Response  string
	UserID    int
	Action    string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ActivityLog) TableName() string {
	return "activity_log"
}

type ActivityLogger struct {
	Id       int    `json:"id"`
	Request  string `json:"request"`
	Response string `json:"Response"`
	UserId   string `json:"user_id"`
}
