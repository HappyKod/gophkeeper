package models

import (
	"time"

	"github.com/google/uuid"
)

type LiteSecret struct {
	ID              uuid.UUID `json:"id"`
	ValueHash       string    `json:"value_hash"`
	DescriptionHash string    `json:"description_hash"`
	IsDeleted       bool      `json:"is_deleted"`
	Ver             time.Time `json:"ver"`
}
