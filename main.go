package main

import (
	"fmt"
	"net/http"
	"project-bioskop/config"
	"project-bioskop/database"
	"project-bioskop/handler"
	"project-bioskop/middleware"
	"project-bioskop/repository"
	"project-bioskop/service"
)

func main() {
	config.LoadConfig()
	database.ConnectDB()
	database.ConnectRedis()

	//database.SeedAll()

	// Repo -> Service -> Handler
	userRepo := repository.NewUserRepository(database.DB)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	trxRepo := repository.NewTransactionRepository(database.DB)
	trxService := service.NewTransactionService(trxRepo)
	trxHandler := handler.NewTransactionHandler(trxService)

    movieRepo := repository.NewMovieRepository(database.DB)
    movieService := service.NewMovieService(movieRepo, database.RDB) 
    movieHandler := handler.NewMovieHandler(movieService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)

	mux.Handle("POST /transactions", middleware.AuthMiddleware(http.HandlerFunc(trxHandler.CreateBooking)))
	mux.Handle("GET /transactions", middleware.AuthMiddleware(http.HandlerFunc(trxHandler.GetHistory)))

	mux.HandleFunc("GET /schedules/{id}/seats", trxHandler.GetSeatsBySchedule)

	mux.HandleFunc("GET /movies", movieHandler.GetMovies)

	fmt.Println("Server running on port :8080")
	http.ListenAndServe(":8080", mux)
}
