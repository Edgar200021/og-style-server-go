package main

import (
	"context"
	cloudinary2 "github.com/cloudinary/cloudinary-go/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net/http"
	"og-style/db"
	"og-style/handlers"
	"og-style/middlewares"
	"og-style/processors"
	"og-style/services"
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

	cloudinary, cldErr := cloudinary2.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if cldErr != nil {
		log.Fatal("Error when trying to connect to cloudinary")
	}

	mux := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "PUT"},
		//AllowedHeaders:             nil,
		//ExposedHeaders:             nil,
		AllowCredentials: true,
		Debug:            true,
	})
	handler := c.Handler(mux)

	var (
		imgUploaderProcessor = services.CldImageUploaderService{Cloudinary: cloudinary}

		userStorage    = db.UserPgStorage{DB: pool}
		cartStorage    = db.CartPgStorage{DB: pool}
		tokenStorage   = db.TokenPgStorage{DB: pool}
		productStorage = db.ProductPgStorage{DB: pool}

		authProcessor    = processors.AuthPgProcessor{&userStorage, &cartStorage, &tokenStorage}
		productProcessor = processors.ProductPgProcessor{&productStorage, &imgUploaderProcessor}

		authHandler    = handlers.AuthHandler{AuthProcessor: &authProcessor}
		productHandler = handlers.ProductHandler{ProductProcessor: &productProcessor}
	)

	mux.HandleFunc("POST /api/v1/auth/sign-up", authHandler.SignUp)
	mux.HandleFunc("POST /api/v1/auth/sign-in", authHandler.SignIn)
	mux.HandleFunc("POST /api/v1/auth/refresh-tokens", authHandler.RefreshTokens)
	mux.HandleFunc("POST /api/v1/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("PATCH /api/v1/auth/reset-password", authHandler.ResetPassword)
	mux.HandleFunc("PATCH /api/v1/auth/update-password", middlewares.Auth(authHandler.UpdatePassword, &userStorage))

	mux.HandleFunc("/api/v1/products", productHandler.GetAll)
	mux.HandleFunc("/api/v1/products/{id}", productHandler.Get)
	mux.HandleFunc("GET /api/v1/products/filters", productHandler.GetFilters)
	mux.HandleFunc("POST /api/v1/products", middlewares.Auth(middlewares.RestrictTo(productHandler.Create, "admin"), &userStorage))
	mux.HandleFunc("PATCH /api/v1/products/{id}", middlewares.Auth(middlewares.RestrictTo(productHandler.Update, "admin"), &userStorage))
	mux.HandleFunc("DELETE /api/v1/products/{id}", middlewares.Auth(middlewares.RestrictTo(productHandler.Delete, "admin"), &userStorage))
	mux.HandleFunc("POST /api/v1/products/upload-image", middlewares.Auth(middlewares.RestrictTo(productHandler.UploadImage, "admin"), &userStorage))

	server := http.Server{
		Addr:        ":4000",
		Handler:     handler,
		ReadTimeout: time.Second * 4,
		//WriteTimeout: time.Second * 5,
		IdleTimeout: time.Second * 120,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Something went wrong")
	}

}
