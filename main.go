package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net/http"
	"og-style/db"
	"og-style/handlers"
	"og-style/processors"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Something went wrong")
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DB_CONNECTION"))
	if err != nil {
		log.Fatal("Error when trying to connect to database")
	}
	defer pool.Close()

	mux := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173/"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "PUT"},
		//AllowedHeaders:             nil,
		//ExposedHeaders:             nil,
		AllowCredentials: true,
		Debug:            true,
	})
	handler := c.Handler(mux)

	var (
		userStorage  = db.UserPgStorage{DB: pool}
		cartStorage  = db.CartPgStorage{DB: pool}
		tokenStorage = db.TokenPgStorage{DB: pool}

		authProcessor = processors.AuthPgProcessor{&userStorage, &cartStorage, &tokenStorage}

		authHandler = handlers.AuthHandler{AuthProcessor: authProcessor}
	)

	mux.HandleFunc("POST /api/v1/auth/sign-up", authHandler.SignUp)
	mux.HandleFunc("POST /api/v1/auth/sign-in", authHandler.SignIn)

	server := http.Server{
		Addr:         ":4000",
		Handler:      handler,
		ReadTimeout:  time.Second * 2,
		WriteTimeout: time.Second * 2,
		IdleTimeout:  time.Second * 120,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Something went wrong")
	}

}
