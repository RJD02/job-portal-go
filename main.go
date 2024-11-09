package main

import (
	"log"
	"net/http"
	"os"

	"RJD02/job-portal/config"
	"RJD02/job-portal/db"
	customMiddleware "RJD02/job-portal/middleware"
	"RJD02/job-portal/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Nevla"))
	// how to import a file in go
}

func setupConfig(config *config.Config) {

	ADMIN_SECRET_KEY := os.Getenv("ADMIN_KEY")

	JWT_SECRET_KEY := os.Getenv("SECRET_KEY")

	client := db.NewClient()

	FROM_GMAIL := os.Getenv("FROM_GMAIL")
	GMAIL_PASSWORD := os.Getenv("GMAIL_PASSWORD")
	TO_GMAIL := os.Getenv("TO_GMAIL")

	ENVIRONMENT := os.Getenv("ENVIRONMENT")

	config.AddSecretKey(JWT_SECRET_KEY)
	config.Connect(client)

	config.AddGmailCreds(FROM_GMAIL, GMAIL_PASSWORD, TO_GMAIL)
	config.SetEnv(ENVIRONMENT)
	config.SetAdminKey(ADMIN_SECRET_KEY)
}

func SetupRouter() *chi.Mux {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading the dot env file")
		return nil
	}

	// app-wide state
	setupConfig(&config.AppConfig)
	r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	r.Use(middleware.Logger)
	r.Use(customMiddleware.CorsMiddleware)
	r.Use(c.Handler)
	r.Get("/", home)

	r.Route("/auth", routes.AuthRouter)
	r.With(customMiddleware.AuthMiddleware).Route("/jobs", routes.JobRouter)

	return r
}

func run() {

	r := SetupRouter()

	log.Println("Current Environment: ", config.AppConfig.ENVIRONMENT)

	log.Println("Connected to database")

	log.Println("Mounted the routes")

	log.Println("Server started on port 5000")
	err := http.ListenAndServe("localhost:5000", r)
	log.Println("I'm failing", err)
}

func main() {
	run()
}
