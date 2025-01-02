package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	// create a router mux
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.Home)

	mux.Post("/authenticate", app.authenticate)
	mux.Post("/signup", app.signup)
	mux.Get("/refresh", app.refreshToken)
	mux.Get("/logout", app.logout)

	// --> deprecated
	mux.Get("/movies", app.AllMovies)
	mux.Get("/movies/{id}", app.GetMovie)
	mux.Get("/genres", app.AllGenres)
	mux.Get("/movies/genres/{id}", app.AllMoviesByGenre)
	mux.Post("/graph", app.moviesGraphQL)
	// -->

	mux.Route("/admin", func(mux chi.Router){
		mux.Use(app.authRequired)

		// --> deprecated
		mux.Get("/movies", app.MovieCatalog)
		mux.Get("/movies/{id}", app.MovieForEdit)
		mux.Put("/movies/0", app.InsertMovie)
		mux.Patch("/movies/{id}", app.UpdateMovie)
		mux.Delete("/movies/{id}", app.DeleteMovie)
		// -->

		// new
		mux.Get("/incomes", app.AllIncomes)
		mux.Post("/incomes/new", app.InsertIncome)
		mux.Get("/sources", app.AllSources)
		mux.Get("/expenses", app.AllExpenses)
		mux.Post("/expenses/new", app.InsertExpense)
		mux.Get("/categories", app.AllCategories)
		mux.Get("/summary", app.GetFinancialSummary)
	})

	return mux
}