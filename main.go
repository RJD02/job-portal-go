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

	JWT_SECRET_KEY := os.Getenv("SECRET_KEY")
	if JWT_SECRET_KEY == "" {
		panic("No SECRET_KEY set")
	}

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatal("Error connecting to database", err)
		panic(err)
	}

	FROM_GMAIL := os.Getenv("FROM_GMAIL")
	GMAIL_PASSWORD := os.Getenv("GMAIL_PASSWORD")
	TO_GMAIL := os.Getenv("TO_GMAIL")

	if FROM_GMAIL == "" || GMAIL_PASSWORD == "" || TO_GMAIL == "" {
		panic("GMAIL Credentials not set")
	}

	ENVIRONMENT := os.Getenv("ENVIRONMENT")

	config.AddSecretKey(JWT_SECRET_KEY)
	config.Connect(client)

	config.AddGmailCreds(FROM_GMAIL, GMAIL_PASSWORD, TO_GMAIL)
	config.SetEnv(ENVIRONMENT)
}

func SetupRouter() *chi.Mux {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading the dot env file")
		return nil
	}

	// db.TestDB(postgres_db)

	// app-wide state
	setupConfig(&config.AppConfig)
	r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	r.Use(middleware.Logger)
	r.Use(c.Handler)
	r.Get("/", home)

	// add admin
	// r.Get("/addadmin", addAdmin)

	r.Route("/auth", routes.AuthRouter)
	r.With(customMiddleware.AuthMiddleware).Route("/jobs", routes.JobRouter)

	return r
}

func run() {
	// init the dotenv

	log.Println("Current Environment: ", config.AppConfig.ENVIRONMENT)

	log.Println("Connected to database")

	r := SetupRouter()
	log.Println("Server started on port 5000")
	http.ListenAndServe("localhost:5000", r)
}

func main() {
	run()
}
