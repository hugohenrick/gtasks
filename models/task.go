package models

import (
	"time"
)

type Task struct {
	ID         uint32     `gorm:"primary_key;auto_increment" json:"id"`
	Title      string     `gorm:"size:200;not null" json:"title"`
	Summary    string     `gorm:"size:2500;not null" json:"summary"`
	UserId     uint32     `gorm:"not null" json:"user_id"`
	Done       bool       `json:"done"`
	User       User       `json:"user,omitempty"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	UpdatedAt  time.Time  `json:"updated_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}
