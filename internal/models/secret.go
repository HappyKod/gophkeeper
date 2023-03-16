package models

import (
	"time"

	"github.com/google/uuid"
)

type Secret struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     string    `json:"owner_id"`
	Value       []byte    `json:"value"`
	Type        string    `json:"secret_type"`
	Description string    `json:"description"`
	IsDeleted   bool      `json:"is_deleted"`
	Ver         time.Time `json:"ver"`
}
