package post

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Data      string
	Timestamp time.Time
	ID        uuid.UUID
}
