package main

import (
	"log"
	"net/http"

	"RJD02/job-portal/config"
	"RJD02/job-portal/db"
	"RJD02/job-portal/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("How are you"))
	// how to import a file in go
}

func addAdmin(w http.ResponseWriter, r *http.Request) {
	// how to import a file in go
	rows, err := config.AppConfig.Db.Query("insert into job_portal.users (username, password) values ('admin', 'admin')")
	if err != nil {
		log.Println("Error in inserting admin", err)
		return
	}
	defer rows.Close()

	w.Write([]byte("Admin added"))
}

func main() {
	// postgres_db init
	postgres_db, err := db.Connect()
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	// db.TestDB(postgres_db)

	// app-wide state
	config.AppConfig.Connect(postgres_db)

	log.Println("Connected to database")

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
	r.Route("/jobs", routes.JobRouter)

	log.Println("Server started on port 5000")
	http.ListenAndServe(":5000", r)
}
