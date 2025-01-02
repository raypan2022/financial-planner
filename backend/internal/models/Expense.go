package models

import "time"

// Category represents a category for expenses.
type Category struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"` // Foreign key to the User table
	Name      string    `json:"name"`    // Name of the category, e.g., "Groceries", "Rent"
	CreatedAt time.Time `json:"-"`       // Timestamp of creation
	UpdatedAt time.Time `json:"-"`       // Timestamp of last update
}

// Expense represents an expense record.
type Expense struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`        // Foreign key to the User table
	Amount        float64   `json:"amount"`         // Amount of the expense
	CategoryID    int       `json:"category_id"`    // Foreign key to the Category table
	Category      *Category `json:"category"`       // Category of the expense
	Date          time.Time `json:"date"`           // Date the expense was incurred
	Description   string    `json:"description"`    // Additional details about the expense
	PaymentMethod string    `json:"payment_method"` // Method of payment, e.g., "Credit Card", "Cash"
	CreatedAt     time.Time `json:"-"`              // Timestamp of creation
	UpdatedAt     time.Time `json:"-"`              // Timestamp of last update
}
