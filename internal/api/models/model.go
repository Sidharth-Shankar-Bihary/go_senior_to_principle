package models

import "time"

// Model is a same type as gorm.Model with addition of json tags to make it easier to generate JSON responses that include its fields.
// Second, User describes simple application user with GORM tags which specify which column the field should be associated with.
type Model struct {
	ID        uint64     `gorm:"primary_key;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
