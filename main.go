package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"RJD02/job-portal/config"
	"RJD02/job-portal/db"
	"RJD02/job-portal/routes"
	"RJD02/job-portal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

//go:generate go run github.com/steebchen/prisma-client-go generate

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

	port := os.Getenv("PORT")

	config.SetPort(port)

	config.AddSecretKey(JWT_SECRET_KEY)
	config.Connect(client)

	config.AddGmailCreds(FROM_GMAIL, GMAIL_PASSWORD, TO_GMAIL)
	config.SetEnv(ENVIRONMENT)
	config.SetAdminKey(ADMIN_SECRET_KEY)
}

func SetupRouter() *chi.Mux {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading the dot env file")
	}

	// app-wide state
	setupConfig(&config.AppConfig)
	r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://jobs-nearme.web.app"},
		AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	r.Use(middleware.Logger)
	r.Use(c.Handler)
	r.Get("/", home)

	r.Route("/auth", routes.AuthRouter)
	r.Route("/jobs", routes.JobRouter)

	return r
}

func run() {

	r := SetupRouter()

	ctx, cancel := context.WithCancel(context.Background())

	go utils.RunUpdateInactiveLinksDaily(ctx)
	defer cancel()

	log.Println("Current Environment: ", config.AppConfig.ENVIRONMENT)

	log.Println("Connected to database")

	log.Println("Mounted the routes")

	log.Println("Server started on port ", config.AppConfig.PORT)
	err := http.ListenAndServe("localhost:"+config.AppConfig.PORT, r)
	log.Println("I'm failing", err)
}

func main() {
	run()
	defer config.AppConfig.Db.Disconnect()
}
