package main

import (
	"backend/internal/graph"
	"backend/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Home displays the status of the api, as JSON.
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go Movies up and running",
		Version: "1.0.0",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// sign up a brand new user and returns a JWT
func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	log.Printf("signup endpoint hit\n")
	var user models.User

	err := app.readJSON(w, r, &user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// validate user against database
	_, err = app.DB.GetUserByEmail(user.Email)
	if err == nil {
			app.errorJSON(w, errors.New("email already exists, try logging in"), http.StatusBadRequest)
			return
	}

	// encrypt password
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
			app.errorJSON(w, errors.New("failed to encrypt password"), http.StatusInternalServerError)
			return
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Password = string(bytes)

	newID, err := app.DB.InsertUser(user)
	if err != nil {
			app.errorJSON(w, err)
			return
	}

	// create a jwt user
	u := jwtUser{
		ID:        newID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	// generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
}

// authenticate authenticates a user when they try to log in, and returns a JWT.
func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	log.Printf("authenticate endpoint hit\n")
	// read json payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate user against database
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// create a jwt user
	u := jwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	// generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
}

// refreshToken checks for a valid refresh cookie, and returns a JWT if it finds one.
func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("refreshToken endpoint hit\n")
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				log.Printf("could not parse claim: %v\n", err)
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// get the user id from the token claims
			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				log.Printf("could not get userid from claim: %v\n", err)
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			user, err := app.DB.GetUserByID(userID)
			if err != nil {
				log.Printf("could not get user %v from db: %v\n", userID, err)
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			u := jwtUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				log.Printf("could not generate tokens: %v\n", err)
				app.errorJSON(w, errors.New("error generating tokens"), http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, app.auth.GetRefreshCookie(tokenPairs.RefreshToken))

			app.writeJSON(w, http.StatusOK, tokenPairs)
		}
	}
}

// logout logs the user out by sending an expired cookie to delete the refresh cookie.
func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	log.Printf("logout endpoint hit\n")
	http.SetCookie(w, app.auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}

// returns all income associated with user
func (app *application) AllIncomes(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
			app.errorJSON(w, errors.New("unable to retrieve claims"), http.StatusUnauthorized)
			return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
			app.errorJSON(w, errors.New("invalid user ID in token"), http.StatusUnauthorized)
			return
	}

	incomes, err := app.DB.AllIncomes(userID)
	if err != nil {
			app.errorJSON(w, err, http.StatusInternalServerError)
			return
	}

	app.writeJSON(w, http.StatusOK, incomes)
}

// returns all expenses associated with user
func (app *application) AllExpenses(w http.ResponseWriter, r *http.Request) {
	log.Printf("AllExpenses endpoint hit\n")
	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
			app.errorJSON(w, errors.New("unable to retrieve claims"), http.StatusUnauthorized)
			return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
			app.errorJSON(w, errors.New("invalid user ID in token"), http.StatusUnauthorized)
			return
	}

	expenses, err := app.DB.AllExpenses(userID)
	if err != nil {
			app.errorJSON(w, err, http.StatusInternalServerError)
			return
	}

	app.writeJSON(w, http.StatusOK, expenses)
}

// insert one paycheque
func (app *application) InsertIncome(w http.ResponseWriter, r *http.Request) {
	log.Printf("InsertIncome endpoint hit\n")
	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
			log.Println("unable to retrieve claims")
			app.errorJSON(w, errors.New("unable to retrieve claims"), http.StatusUnauthorized)
			return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
			log.Println("invalid user ID in token")
			app.errorJSON(w, errors.New("invalid user ID in token"), http.StatusUnauthorized)
			return
	}

	var income models.Income
	err = app.readJSON(w, r, &income)
	if err != nil {
			log.Printf("error reading JSON: %s\n", err)
			app.errorJSON(w, err)
			return
	}

	income.UserID = userID
	income.CreatedAt = time.Now()
	income.UpdatedAt = time.Now()
	income.Source.UserID = userID
	income.Source.CreatedAt = time.Now()
	income.Source.UpdatedAt = time.Now()

	err = app.DB.InsertIncome(&income)
	if err != nil {
			log.Println("error inserting income")
			app.errorJSON(w, err)
			return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "income inserted",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

// insert one expense
func (app *application) InsertExpense(w http.ResponseWriter, r *http.Request) {
	log.Printf("InsertExpense endpoint hit\n")
	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
			app.errorJSON(w, errors.New("unable to retrieve claims"), http.StatusUnauthorized)
			return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
			app.errorJSON(w, errors.New("invalid user ID in token"), http.StatusUnauthorized)
			return
	}

	var expense models.Expense
	err = app.readJSON(w, r, &expense)
	if err != nil {
			app.errorJSON(w, err)
			return
	}

	expense.UserID = userID
	expense.CreatedAt = time.Now()
	expense.UpdatedAt = time.Now()
	expense.Category.UserID = userID
	expense.Category.CreatedAt = time.Now()
	expense.Category.UpdatedAt = time.Now()

	err = app.DB.InsertExpense(&expense)
	if err != nil {
			app.errorJSON(w, err)
			return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "expense inserted",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

// get all sources belonging to user
func (app *application) AllSources(w http.ResponseWriter, r *http.Request) {
	log.Printf("AllSources endpoint hit\n")
	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
			app.errorJSON(w, errors.New("unable to retrieve claims"), http.StatusUnauthorized)
			return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
			app.errorJSON(w, errors.New("invalid user ID in token"), http.StatusUnauthorized)
			return
	}

	sources, err := app.DB.AllSources(userID)
	if err != nil {
			app.errorJSON(w, err)
			return
	}

	app.writeJSON(w, http.StatusOK, sources)
}

// get all categories belonging to user
func (app *application) AllCategories(w http.ResponseWriter, r *http.Request) {
	log.Printf("AllCategories endpoint hit\n")
	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
			app.errorJSON(w, errors.New("unable to retrieve claims"), http.StatusUnauthorized)
			return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
			app.errorJSON(w, errors.New("invalid user ID in token"), http.StatusUnauthorized)
			return
	}

	categories, err := app.DB.AllCategories(userID)
	if err != nil {
			app.errorJSON(w, err)
			return
	}

	app.writeJSON(w, http.StatusOK, categories)
}

// get summary for dashboard
func (app *application) GetFinancialSummary(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetFinancialSummary endpoint hit\n")
	// get the userId using the jwt token
	userID, err := app.getUserIDFromContext(r)
	if err != nil {
			app.errorJSON(w, err)
			return
	}

	// Get overall sums
	totalIncome, err := app.DB.GetTotalIncome(userID)
	if err != nil {
			app.errorJSON(w, err)
			return
	}

	totalExpenses, err := app.DB.GetTotalExpenses(userID)
	if err != nil {
			app.errorJSON(w, err)
			return
	}

	accountBalance := totalIncome - totalExpenses

	// Get income by source and expenses by category
	incomeBySource, err := app.DB.GetIncomeBySource(userID)
	if err != nil {
			app.errorJSON(w, err)
			return
	}

	expensesByCategory, err := app.DB.GetExpensesByCategory(userID)
	if err != nil {
			app.errorJSON(w, err)
			return
	}

	// Get data for the past 12 months
	var months []map[string]interface{}
	for i := 11; i >= 0; i-- {
			monthData := make(map[string]interface{})
			
			monthName := time.Now().AddDate(0, -i, 0).Format("January 2006")
			monthData["month"] = monthName

			incomeThisMonth, err := app.DB.GetIncomeForMonth(userID, i)
			if err != nil {
					app.errorJSON(w, err)
					return
			}

			expensesThisMonth, err := app.DB.GetExpensesForMonth(userID, i)
			if err != nil {
					app.errorJSON(w, err)
					return
			}

			netIncomeThisMonth := incomeThisMonth - expensesThisMonth

			incomeBySourceThisMonth, err := app.DB.GetIncomeBySourceForMonth(userID, i)
			if err != nil {
					app.errorJSON(w, err)
					return
			}

			expensesByCategoryThisMonth, err := app.DB.GetExpensesByCategoryForMonth(userID, i)
			if err != nil {
					app.errorJSON(w, err)
					return
			}

			top3IncomeSources, err := app.DB.GetTop3IncomeSourcesForMonth(userID, i)
			if err != nil {
					app.errorJSON(w, err)
					return
			}

			top3ExpenseCategories, err := app.DB.GetTop3ExpenseCategoriesForMonth(userID, i)
			if err != nil {
					app.errorJSON(w, err)
					return
			}

			monthData["net_income"] = netIncomeThisMonth
			monthData["income_sum"] = incomeThisMonth
			monthData["expense_sum"] = expensesThisMonth
			monthData["income_by_source"] = incomeBySourceThisMonth
			monthData["expense_by_category"] = expensesByCategoryThisMonth
			monthData["top3IncomeThisMonth"] = top3IncomeSources
			monthData["top3ExpenseThisMonth"] = top3ExpenseCategories

			months = append(months, monthData)
	}

	// Build the final JSON response
	summary := map[string]interface{}{
			"account_balance":       accountBalance,
			"income_sum_total":      totalIncome,
			"expense_sum_total":     totalExpenses,
			"overall_income_by_source": incomeBySource,
			"overall_expense_by_category": expensesByCategory,
			"months":                months,
	}

	app.writeJSON(w, http.StatusOK, summary)
}



// -------------------------------------- OLD CODE FOR REFERENCE --------------------------------------


// AllMovies returns a slice of all movies as JSON.
func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.AllMovies()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movies)
}

// MovieCatalog returns a list of all movies as JSON
func (app *application) MovieCatalog(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.AllMovies()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movies)
}

// GetMovie returns one movie, as JSON.
func (app *application) GetMovie(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie, err := app.DB.OneMovie(movieID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movie)
}

// MovieForEdit returns a JSON payload for a given movie and a list of all genres, for edit.
func (app *application) MovieForEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie, genres, err := app.DB.OneMovieForEdit(movieID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload = struct {
		Movie  *models.Movie   `json:"movie"`
		Genres []*models.Genre `json:"genres"`
	}{
		movie,
		genres,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// AllGenres returns a slice of all genres as JSON.
func (app *application) AllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.DB.AllGenres()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, genres)
}

// InsertMovie receives a JSON payload and tries to insert a movie into the database.
func (app *application) InsertMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie

	err := app.readJSON(w, r, &movie)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// try to get an image
	movie = app.getPoster(movie)

	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()

	newID, err := app.DB.InsertMovie(movie)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// now handle genres
	err = app.DB.UpdateMovieGenres(newID, movie.GenresArray)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "movie updated",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

// getPoster tries to get a poster image from themoviedb.org.
func (app *application) getPoster(movie models.Movie) models.Movie {
	type TheMovieDB struct {
		Page    int `json:"page"`
		Results []struct {
			PosterPath string `json:"poster_path"`
		} `json:"results"`
		TotalPages int `json:"total_pages"`
	}

	client := &http.Client{}
	theUrl := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s", app.APIKey)

	// https://api.themoviedb.org/3/search/movie?api_key=b41447e6319d1cd467306735632ba733&query=Die+Hard

	req, err := http.NewRequest("GET", theUrl+"&query="+url.QueryEscape(movie.Title), nil)
	if err != nil {
		log.Println(err)
		return movie
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return movie
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return movie
	}

	var responseObject TheMovieDB

	json.Unmarshal(bodyBytes, &responseObject)

	if len(responseObject.Results) > 0 {
		movie.Image = responseObject.Results[0].PosterPath
	}

	return movie
}

// UpdateMovie updates a movie in the database, based on a JSON payload.
func (app *application) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	var payload models.Movie

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie, err := app.DB.OneMovie(payload.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie.Title = payload.Title
	movie.ReleaseDate = payload.ReleaseDate
	movie.Description = payload.Description
	movie.MPAARating = payload.MPAARating
	movie.RunTime = payload.RunTime
	movie.UpdatedAt = time.Now()

	err = app.DB.UpdateMovie(*movie)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.UpdateMovieGenres(movie.ID, payload.GenresArray)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "movie updated",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

// DeleteMovie deletes a movie from the database, by ID.
func (app *application) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.DeleteMovie(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "movie deleted",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) AllMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movies, err := app.DB.AllMovies(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, movies)
}


// -------------------------------------- GRAPHQL --------------------------------------


func (app *application) moviesGraphQL(w http.ResponseWriter, r *http.Request) {
	// we need to populate our Graph type with the movies
	movies, _ := app.DB.AllMovies()

	// get the query from the request
	q, _ := io.ReadAll(r.Body)
	query := string(q)

	// create a new variable of type *graph.Graph
	g := graph.New(movies)

	// set the query string on the variable
	g.QueryString = query

	// perform the query
	resp, err := g.Query()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// send the response
	j, _ := json.MarshalIndent(resp, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}