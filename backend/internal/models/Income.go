package models

import "time"

// Source represents a source of income.
type Source struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"` // Foreign key to the User table
	Name      string    `json:"name"`    // Name of the source, e.g., "Salary", "Freelance"
	CreatedAt time.Time `json:"-"`       // Timestamp of creation
	UpdatedAt time.Time `json:"-"`       // Timestamp of last update
}

// Income represents an income record.
type Income struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`     // Foreign key to the User table
	Amount      float64   `json:"amount"`      // Amount of the income
	SourceID    int       `json:"source_id"`   // ID of the source
	Source      *Source   `json:"source"`      // Source of income
	Date        time.Time `json:"date"`        // Date the income was received
	Description string    `json:"description"` // Additional details about the income
	CreatedAt   time.Time `json:"-"`           // Timestamp of creation
	UpdatedAt   time.Time `json:"-"`           // Timestamp of last update
}
