package main

import (
	"cinerdle-back/internal/adapters/email"
	"cinerdle-back/internal/adapters/repository"
	adapters "cinerdle-back/internal/adapters/tmdb"
	delivery "cinerdle-back/internal/delivery/http"
	"cinerdle-back/internal/usecases"
	"cinerdle-back/pkg/config"
	"cinerdle-back/pkg/jwt"
	"cinerdle-back/pkg/password"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	cfg := config.LoadConfig()
	db, err := sql.Open("postgres", cfg.DatabaseURL)

	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	userRepo := repository.NewPostgresUserRepository(db)
	emailRepo := repository.NewPostgreEmailVerificationRepository(db)
	passwordHasher := password.NewBcryptHasher(14)
	jwtManager := jwt.NewJwtManagerImpl(cfg.JWTSecret)
	emailVerification := email.NewSmtpEmailSender(
		cfg.SMTPAddress, cfg.SMTPPort, cfg.SMTPLogin, cfg.SMTPPassword,
	)
	authCase := usecases.NewAuthCase(
		userRepo, emailRepo, passwordHasher, jwtManager, emailVerification,
	)
	authHand := delivery.NewAuthHandler(authCase)
	authMiddleCase := delivery.NewAuthMiddleware(jwtManager)
	tmdbClient := adapters.NewTMDbClient(
		cfg.TMDbAPIKey,
		http.Client{},
	)

	movieRepo := repository.NewPostgresMovieRepository(db)
	movieCase := usecases.NewMovieUseCase(
		tmdbClient, movieRepo,
	)
	movieHand := delivery.NewMovieHandler(movieCase)

	router := delivery.NewRouter(authHand, authMiddleCase, movieHand)

	mux := http.NewServeMux()
	router.Setup(mux)
	fmt.Println("Server started on :8080")
	http.ListenAndServe(cfg.Port, mux)

}
