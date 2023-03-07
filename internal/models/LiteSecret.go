package models

import "time"

type LiteSecret struct {
	ID              int       `json:"id"`
	ValueHash       string    `json:"value_hash"`
	DescriptionHash string    `json:"description_hash"`
	IsDeleted       bool      `json:"is_deleted"`
	Ver             time.Time `json:"ver"`
}
