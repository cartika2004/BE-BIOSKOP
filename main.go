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
	// 1. Init
	config.LoadConfig()
	database.ConnectDB()
	database.ConnectRedis()

	// Hapus atau comment baris ini kalau data sudah ada biar gak seeding ulang terus
	//database.SeedAll()

	// 2. Dependency Injection (Wiring)
	// Repo -> Service -> Handler
	userRepo := repository.NewUserRepository(database.DB)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// --- SETUP TRANSACTION (YANG BARU) ---
	trxRepo := repository.NewTransactionRepository(database.DB)
	trxService := service.NewTransactionService(trxRepo)
	trxHandler := handler.NewTransactionHandler(trxService)

	// --- SETUP MOVIE (BARU) ---
    movieRepo := repository.NewMovieRepository(database.DB)
    movieService := service.NewMovieService(movieRepo, database.RDB) // Inject Redis di sini!
    movieHandler := handler.NewMovieHandler(movieService)

	// 3. Routing Manual (Tanpa library tambahan seperti Mux/Chi)
	mux := http.NewServeMux()

	// Rute Publik (Gak butuh token)
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)

	// --- RUTE BOOKING (DIPROTEKSI) ---
	// Artinya: Hanya user yang punya token yang bisa akses ini
	mux.Handle("POST /transactions", middleware.AuthMiddleware(http.HandlerFunc(trxHandler.CreateBooking)))
	mux.Handle("GET /transactions", middleware.AuthMiddleware(http.HandlerFunc(trxHandler.GetHistory)))

	mux.HandleFunc("GET /schedules/{id}/seats", trxHandler.GetSeatsBySchedule)

	mux.HandleFunc("GET /movies", movieHandler.GetMovies)

	fmt.Println("🚀 Server running on port :8080")
	http.ListenAndServe(":8080", mux)
}
