package domain

import (
	"time"
)

type Post struct {
	Data      string
	Timestamp time.Time
	ID        uint `gorm:"primaryKey;autoIncrement"`
}
