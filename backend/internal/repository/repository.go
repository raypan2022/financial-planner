package repository

import (
	"backend/internal/models"
	"database/sql"
)


type DatabaseRepo interface {
	Connection() *sql.DB
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	AllIncomes(id int) ([]*models.Income, error)
	AllExpenses(id int) ([]*models.Expense, error)
	InsertIncome(income *models.Income) error
	InsertExpense(expense *models.Expense) error
	AllSources(id int) ([]*models.Source, error)
	AllCategories(id int) ([]*models.Category, error)
	GetTotalIncome(userID int) (float64, error)
	GetTotalExpenses(userID int) (float64, error)
	GetIncomeBySource(userID int) (map[string]float64, error)
	GetExpensesByCategory(userID int) (map[string]float64, error)
	GetIncomeForMonth(userID, monthsAgo int) (float64, error)
	GetExpensesForMonth(userID, monthsAgo int) (float64, error)
	GetIncomeBySourceForMonth(userID, monthsAgo int) (map[string]float64, error)
	GetExpensesByCategoryForMonth(userID, monthsAgo int) (map[string]float64, error)
	GetTop3IncomeSourcesForMonth(userID, monthsAgo int) ([]map[string]interface{}, error)
	GetTop3ExpenseCategoriesForMonth(userID, monthsAgo int) ([]map[string]interface{}, error)

	// ----------------- NEPRECATED OLD CODE -----------------

	AllMovies(genre ...int) ([]*models.Movie, error)
	OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error)
	OneMovie(id int) (*models.Movie, error)
	AllGenres() ([]*models.Genre, error)
	InsertMovie(movie models.Movie) (int, error)
	InsertUser(user models.User) (int, error)
	UpdateMovieGenres(id int, genreIDs []int) error
	UpdateMovie(movie models.Movie) error
	DeleteMovie(id int) error
}