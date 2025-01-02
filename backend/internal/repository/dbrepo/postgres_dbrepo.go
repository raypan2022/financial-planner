package dbrepo

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// PostgresDBRepo is the struct used to wrap our database connection pool, so that we
// can easily swap out a real database for a test database, or move to another database
// entirely, as long as the thing being swapped implements all of the functions in the type
// repository.DatabaseRepo.
type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

// Connection returns underlying connection pool.
func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

// GetUserByEmail returns one use, by email.
func (m *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password,
			created_at, updated_at from users where email = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID returns one use, by ID.
func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password,
			created_at, updated_at from users where id = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *PostgresDBRepo) InsertUser(user models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into users (first_name, last_name, email, password, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6) returning id`
	
	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&newID)

	if err != nil {
		log.Printf("Error inserting user: %v\n", err)
		return 0, err
	}

	return newID, nil
}

func (m *PostgresDBRepo) AllIncomes(id int) ([]*models.Income, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, user_id, amount, source_id, date, description, created_at, updated_at from incomes where user_id = $1 order by date desc`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []*models.Income

	for rows.Next() {
		var income models.Income
		err := rows.Scan(
			&income.ID,
			&income.UserID,
			&income.Amount,
			&income.SourceID,
			&income.Date,
			&income.Description,
			&income.CreatedAt,
			&income.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Fetch the corresponding source record for this income
		sourceQuery := `select id, name, created_at, updated_at from sources where id = $1`
		var source models.Source
		err = m.DB.QueryRowContext(ctx, sourceQuery, income.SourceID).Scan(
			&source.ID,
			&source.Name,
			&source.CreatedAt,
			&source.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set the Source property of the income
		income.Source = &source

		incomes = append(incomes, &income)
	}

	return incomes, nil
}

func (m *PostgresDBRepo) AllExpenses(id int) ([]*models.Expense, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, user_id, amount, category_id, date, description, payment_method, created_at, updated_at from expenses where user_id = $1 order by date desc`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []*models.Expense

	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(
			&expense.ID,
			&expense.UserID,
			&expense.Amount,
			&expense.CategoryID,
			&expense.Date,
			&expense.Description,
			&expense.PaymentMethod,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Fetch the corresponding category record for this expense
		categoryQuery := `select id, name, created_at, updated_at from categories where id = $1`
		var category models.Category
		err = m.DB.QueryRowContext(ctx, categoryQuery, expense.CategoryID).Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set the Category property of the expense
		expense.Category = &category

		expenses = append(expenses, &expense)
	}

	return expenses, nil
}

func (m *PostgresDBRepo) InsertIncome(income *models.Income) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Check if the source exists for the specific user, and insert if it doesn't
	var sourceID int

	err := m.DB.QueryRowContext(ctx, `
		SELECT id FROM sources WHERE name = $1 AND user_id = $2`, 
		income.Source.Name, income.UserID).Scan(&sourceID)
			
	if err == sql.ErrNoRows {
		// Source doesn't exist for this user, insert it
		err = m.DB.QueryRowContext(ctx, `
			INSERT INTO sources (name, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`, 
			income.Source.Name, income.UserID, income.Source.CreatedAt, income.Source.UpdatedAt).Scan(&sourceID)
		if err != nil {
			log.Printf("created at: %v\n", income.Source.CreatedAt)
			log.Printf("Error inserting source: %v\n", err)
			return err
		}
	} else if err != nil {
		log.Printf("Error checking if source exists: %v\n", err)
		return err
	}

	income.SourceID = sourceID

	// Insert the income record
	query := `
			INSERT INTO incomes (user_id, amount, source_id, date, description, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = m.DB.ExecContext(ctx, query, income.UserID, income.Amount, income.SourceID, income.Date, income.Description, income.CreatedAt, income.UpdatedAt)
	if err != nil {
		log.Printf("Error inserting income: %v\n", err)
	}
	return err
}

func (m *PostgresDBRepo) InsertExpense(expense *models.Expense) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Check if the category exists for the specific user, and insert if it doesn't
	var categoryID int
	err := m.DB.QueryRowContext(ctx, `
			SELECT id FROM categories WHERE name = $1 AND user_id = $2`, 
			expense.Category.Name, expense.UserID).Scan(&categoryID)
	if err == sql.ErrNoRows {
			// Category doesn't exist for this user, insert it
			err = m.DB.QueryRowContext(ctx, `
					INSERT INTO categories (name, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`, 
					expense.Category.Name, expense.UserID, expense.Category.CreatedAt, expense.Category.UpdatedAt).Scan(&categoryID)
			if err != nil {
					return err
			}
	} else if err != nil {
			return err
	}
	expense.CategoryID = categoryID

	// Insert the expense record
	query := `
			INSERT INTO expenses (user_id, amount, category_id, date, description, payment_method, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = m.DB.ExecContext(ctx, query, expense.UserID, expense.Amount, expense.CategoryID, expense.Date, expense.Description, expense.PaymentMethod, expense.CreatedAt, expense.UpdatedAt)
	return err
}

func (m *PostgresDBRepo) AllSources(id int) ([]*models.Source, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, created_at, updated_at from sources where user_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []*models.Source

	for rows.Next() {
		var source models.Source
		err := rows.Scan(
			&source.ID,
			&source.Name,
			&source.CreatedAt,
			&source.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sources = append(sources, &source)
	}

	return sources, nil
}

func (m *PostgresDBRepo) AllCategories(id int) ([]*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, created_at, updated_at from categories where user_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category

	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (m *PostgresDBRepo) GetTotalIncome(userID int) (float64, error) {
	var totalIncome float64
	query := `SELECT COALESCE(SUM(amount), 0) FROM incomes WHERE user_id = $1`
	
	err := m.DB.QueryRow(query, userID).Scan(&totalIncome)
	if err != nil {
			return 0, err
	}
	
	return totalIncome, nil
}

func (m *PostgresDBRepo) GetTotalExpenses(userID int) (float64, error) {
	var totalExpenses float64
	query := `SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE user_id = $1`
	
	err := m.DB.QueryRow(query, userID).Scan(&totalExpenses)
	if err != nil {
			return 0, err
	}
	
	return totalExpenses, nil
}

func (m *PostgresDBRepo) GetIncomeBySource(userID int) (map[string]float64, error) {
	query := `SELECT s.name, COALESCE(SUM(i.amount), 0) FROM incomes i
						JOIN sources s ON i.source_id = s.id
						WHERE i.user_id = $1 GROUP BY s.name`
	
	rows, err := m.DB.Query(query, userID)
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	incomeBySource := make(map[string]float64)
	
	for rows.Next() {
			var sourceName string
			var amount float64
			
			err := rows.Scan(&sourceName, &amount)
			if err != nil {
					return nil, err
			}
			
			incomeBySource[sourceName] = amount
	}
	
	return incomeBySource, nil
}

func (m *PostgresDBRepo) GetExpensesByCategory(userID int) (map[string]float64, error) {
	query := `SELECT c.name, COALESCE(SUM(e.amount), 0) FROM expenses e
						JOIN categories c ON e.category_id = c.id
						WHERE e.user_id = $1 GROUP BY c.name`
	
	rows, err := m.DB.Query(query, userID)
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	expensesByCategory := make(map[string]float64)
	
	for rows.Next() {
			var categoryName string
			var amount float64
			
			err := rows.Scan(&categoryName, &amount)
			if err != nil {
					return nil, err
			}
			
			expensesByCategory[categoryName] = amount
	}
	
	return expensesByCategory, nil
}

func (m *PostgresDBRepo) GetIncomeForMonth(userID, monthsAgo int) (float64, error) {
	var income float64
	query := `SELECT COALESCE(SUM(amount), 0) FROM incomes
              WHERE user_id = $1 AND date_trunc('month', date) = date_trunc('month', (CURRENT_DATE - INTERVAL '1 month' * $2))`
	
	err := m.DB.QueryRow(query, userID, monthsAgo).Scan(&income)
	if err != nil {
			log.Println(err)
			return 0, err
	}
	
	return income, nil
}

func (m *PostgresDBRepo) GetExpensesForMonth(userID, monthsAgo int) (float64, error) {
	var expenses float64
	query := `SELECT COALESCE(SUM(amount), 0) FROM expenses
						WHERE user_id = $1 AND date_trunc('month', date) = date_trunc('month', (CURRENT_DATE - INTERVAL '1 month' * $2))`
	
	
	err := m.DB.QueryRow(query, userID, monthsAgo).Scan(&expenses)
	if err != nil {
			return 0, err
	}
	
	return expenses, nil
}

func (m *PostgresDBRepo) GetIncomeBySourceForMonth(userID, monthsAgo int) (map[string]float64, error) {
	query := `SELECT s.name, COALESCE(SUM(i.amount), 0) FROM incomes i
						JOIN sources s ON i.source_id = s.id
						WHERE i.user_id = $1 AND date_trunc('month', i.date) = date_trunc('month', (CURRENT_DATE - INTERVAL '1 month' * $2))
						GROUP BY s.name`
	
	rows, err := m.DB.Query(query, userID, monthsAgo)
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	incomeBySource := make(map[string]float64)
	
	for rows.Next() {
			var sourceName string
			var amount float64
			
			err := rows.Scan(&sourceName, &amount)
			if err != nil {
					return nil, err
			}
			
			incomeBySource[sourceName] = amount
	}
	
	return incomeBySource, nil
}

func (m *PostgresDBRepo) GetExpensesByCategoryForMonth(userID, monthsAgo int) (map[string]float64, error) {
	query := `SELECT c.name, COALESCE(SUM(e.amount), 0) FROM expenses e
						JOIN categories c ON e.category_id = c.id
						WHERE e.user_id = $1 AND date_trunc('month', e.date) = date_trunc('month', (CURRENT_DATE - INTERVAL '1 month' * $2))
						GROUP BY c.name`
	
	rows, err := m.DB.Query(query, userID, monthsAgo)
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	expensesByCategory := make(map[string]float64)
	
	for rows.Next() {
			var categoryName string
			var amount float64
			
			err := rows.Scan(&categoryName, &amount)
			if err != nil {
					return nil, err
			}
			
			expensesByCategory[categoryName] = amount
	}
	
	return expensesByCategory, nil
}

func (m *PostgresDBRepo) GetTop3IncomeSourcesForMonth(userID, monthsAgo int) ([]map[string]interface{}, error) {
	query := `SELECT s.name, COALESCE(SUM(i.amount), 0) FROM incomes i
						JOIN sources s ON i.source_id = s.id
						WHERE i.user_id = $1 AND date_trunc('month', i.date) = date_trunc('month', (CURRENT_DATE - INTERVAL '1 month' * $2))
						GROUP BY s.name ORDER BY SUM(i.amount) DESC LIMIT 3`
	
	rows, err := m.DB.Query(query, userID, monthsAgo)
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	var top3IncomeSources []map[string]interface{}
	
	for rows.Next() {
			var sourceName string
			var amount float64
			
			err := rows.Scan(&sourceName, &amount)
			if err != nil {
					return nil, err
			}
			
			top3IncomeSources = append(top3IncomeSources, map[string]interface{}{
					"source": sourceName,
					"amount": amount,
			})
	}
	
	return top3IncomeSources, nil
}

func (m *PostgresDBRepo) GetTop3ExpenseCategoriesForMonth(userID, monthsAgo int) ([]map[string]interface{}, error) {
	query := `SELECT c.name, COALESCE(SUM(e.amount), 0) FROM expenses e
						JOIN categories c ON e.category_id = c.id
						WHERE e.user_id = $1 AND date_trunc('month', e.date) = date_trunc('month', (CURRENT_DATE - INTERVAL '1 month' * $2))
						GROUP BY c.name ORDER BY SUM(e.amount) DESC LIMIT 3`
	
	rows, err := m.DB.Query(query, userID, monthsAgo)
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	var top3ExpenseCategories []map[string]interface{}
	
	for rows.Next() {
			var categoryName string
			var amount float64
			
			err := rows.Scan(&categoryName, &amount)
			if err != nil {
					return nil, err
			}
			
			top3ExpenseCategories = append(top3ExpenseCategories, map[string]interface{}{
					"category": categoryName,
					"amount": amount,
			})
	}
	
	return top3ExpenseCategories, nil
}


// -------------------------------------- OLD CODE FOR REFERENCE --------------------------------------


// AllMovies returns a slice of movies, sorted by name. If the optional parameter genre
// is supplied, then only all movies for a particular genre is returned.
func (m *PostgresDBRepo) AllMovies(genre ...int) ([]*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	where := ""
	if len(genre) > 0 {
		where = fmt.Sprintf("where id in (select movie_id from movies_genres where genre_id = %d)", genre[0])
	}

	query := fmt.Sprintf(`
		select
			id, title, release_date, runtime,
			mpaa_rating, description, coalesce(image, ''),
			created_at, updated_at
		from
			movies %s
		order by
			title
	`, where)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*models.Movie

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.RunTime,
			&movie.MPAARating,
			&movie.Description,
			&movie.Image,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	return movies, nil
}

// OneMovie returns a single movie and associated genres, if any.
func (m *PostgresDBRepo) OneMovie(id int) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, title, release_date, runtime, mpaa_rating, 
		description, coalesce(image, ''), created_at, updated_at
		from movies where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var movie models.Movie

	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.RunTime,
		&movie.MPAARating,
		&movie.Description,
		&movie.Image,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// get genres, if any
	query = `select g.id, g.genre from movies_genres mg
		left join genres g on (mg.genre_id = g.id)
		where mg.movie_id = $1
		order by g.genre`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	var genres []*models.Genre
	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, err
		}

		genres = append(genres, &g)
	}

	movie.Genres = genres

	return &movie, err
}

// OneMovieForEdit returns a single movie and associated genres, if any, for edit.
func (m *PostgresDBRepo) OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, title, release_date, runtime, mpaa_rating, 
		description, coalesce(image, ''), created_at, updated_at
		from movies where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var movie models.Movie

	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.RunTime,
		&movie.MPAARating,
		&movie.Description,
		&movie.Image,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		return nil, nil, err
	}

	// get genres, if any
	query = `select g.id, g.genre from movies_genres mg
		left join genres g on (mg.genre_id = g.id)
		where mg.movie_id = $1
		order by g.genre`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, nil, err
	}
	defer rows.Close()

	var genres []*models.Genre
	var genresArray []int

	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, nil, err
		}

		genres = append(genres, &g)
		genresArray = append(genresArray, g.ID)
	}

	movie.Genres = genres
	movie.GenresArray = genresArray

	var allGenres []*models.Genre

	query = "select id, genre from genres order by genre"
	gRows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	defer gRows.Close()

	for gRows.Next() {
		var g models.Genre
		err := gRows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, nil, err
		}

		allGenres = append(allGenres, &g)
	}

	return &movie, allGenres, err
}

// AllGenres returns a slice of genres, sorted by name.
func (m *PostgresDBRepo) AllGenres() ([]*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, genre, created_at, updated_at from genres order by genre`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []*models.Genre

	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		genres = append(genres, &g)
	}

	return genres, nil
}

// InsertMovie inserts one movie into the database.
func (m *PostgresDBRepo) InsertMovie(movie models.Movie) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into movies (title, description, release_date, runtime,
			mpaa_rating, created_at, updated_at, image)
			values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		movie.Title,
		movie.Description,
		movie.ReleaseDate,
		movie.RunTime,
		movie.MPAARating,
		movie.CreatedAt,
		movie.UpdatedAt,
		movie.Image,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// UpdateMovie updates one movie in the database.
func (m *PostgresDBRepo) UpdateMovie(movie models.Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update movies set title = $1, description = $2, release_date = $3,
				runtime = $4, mpaa_rating = $5,
				updated_at = $6, image = $7 where id = $8`

	_, err := m.DB.ExecContext(ctx, stmt,
		movie.Title,
		movie.Description,
		movie.ReleaseDate,
		movie.RunTime,
		movie.MPAARating,
		movie.UpdatedAt,
		movie.Image,
		movie.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// UpdateMovieGenres first deletes all genres associated with a movie, and
// then inserts the ones stored in genreIDs.
func (m *PostgresDBRepo) UpdateMovieGenres(id int, genreIDs []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from movies_genres where movie_id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	for _, n := range genreIDs {
		stmt := `insert into movies_genres (movie_id, genre_id) values ($1, $2)`
		_, err := m.DB.ExecContext(ctx, stmt, id, n)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMovie deletes one movie, by id.
func (m *PostgresDBRepo) DeleteMovie(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from movies where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}
