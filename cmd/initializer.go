package main

import (
	_ "context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/google/uuid"
	_ "github.com/joho/godotenv"
	_ "google.golang.org/api/option"
	"log"
	"net/http"

	"workout/internal/handlers"
	_ "workout/internal/handlers"
	_ "workout/internal/models"
	"workout/internal/repositories"
	_ "workout/internal/repositories"
	"workout/internal/services"
	_ "workout/internal/services"
	_ "workout/utils"
)

type application struct {

	errorLog         *log.Logger
	infoLog          *log.Logger
	userHandler      *handlers.UserHandler
	userRepo         *repositories.UserRepository
	programHandler   *handlers.ProgramHandler
	programRepo      *repositories.ProgramRepository
	dayHandler       *handlers.DayHandler
	dayRepo          *repositories.DayRepository
	exerciseHandler  *handlers.ExerciseHandler
	exerciseRepo     *repositories.ExerciseRepository
	foodHandler      *handlers.FoodHandler
	foodRepo         *repositories.FoodRepository
	inviteHandler    *handlers.InviteHandler
	inviteRepo       *repositories.InviteRepository
	analyticsHandler *handlers.AnalyticsHandler
	analyticsRepo    *repositories.AnalyticsRepository

}

func initializeApp(db *sql.DB, errorLog, infoLog *log.Logger) *application {
	// Repositories
	userRepo := repositories.UserRepository{DB: db}
	programRepo := repositories.ProgramRepository{DB: db}
	dayRepo := repositories.DayRepository{DB: db}
	exerciseRepo := repositories.ExerciseRepository{DB: db}
	foodRepo := repositories.FoodRepository{DB: db}
	inviteRepo := repositories.InviteRepository{DB: db}

	analyticsRepo := repositories.AnalyticsRepository{DB: db}


	// Services
	userService := &services.UserService{UserRepo: &userRepo}
	programService := &services.ProgramService{Repo: &programRepo}
	dayService := &services.DayService{Repo: &dayRepo}
	exerciseService := &services.ExerciseService{Repo: &exerciseRepo}
	foodService := &services.FoodService{Repo: &foodRepo}
	inviteService := &services.InviteService{Repo: &inviteRepo, UserRepo: &userRepo}

	analyticsService := &services.AnalyticsService{Repo: &analyticsRepo}


	// Handlers
	userHandler := &handlers.UserHandler{Service: userService}
	programHandler := &handlers.ProgramHandler{Service: programService}
	dayHandler := &handlers.DayHandler{Service: dayService}
	exerciseHandler := &handlers.ExerciseHandler{Service: exerciseService}
	foodHandler := &handlers.FoodHandler{Service: foodService}
	inviteHandler := &handlers.InviteHandler{Service: inviteService}
	analyticsHandler := &handlers.AnalyticsHandler{Service: analyticsService}

	return &application{
		errorLog:         errorLog,
		infoLog:          infoLog,
		userHandler:      userHandler,
		programHandler:   programHandler,
		userRepo:         &userRepo,
		programRepo:      &programRepo,
		dayRepo:          &dayRepo,
		dayHandler:       dayHandler,
		exerciseRepo:     &exerciseRepo,
		exerciseHandler:  exerciseHandler,
		foodRepo:         &foodRepo,
		foodHandler:      foodHandler,
		inviteRepo:       &inviteRepo,
		inviteHandler:    inviteHandler,
		analyticsRepo:    &analyticsRepo,
		analyticsHandler: analyticsHandler,
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Failed to open DB: %v", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Printf("Failed to ping DB: %v", err)
		return nil, err
	}
	db.SetMaxIdleConns(35)
	log.Println("Successfully connected to database")
	return db, nil
}

func addSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}
