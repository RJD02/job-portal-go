package main

import (
	"log"
	"net/http"
	"os"

	"RJD02/job-portal/config"
	"RJD02/job-portal/db"
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

//	func addAdmin(w http.ResponseWriter, r *http.Request) {
//		ctx := context.Background()
//		// how to import a file in go
//		createdUser, err := config.AppConfig.Db.User.CreateOne(
//			db.User.Username.Set("admin"),
//			db.User.Password.Set("admin"),
//			db.User.Email.Set("admin@admin.com"),
//		).Exec(ctx)
//		// rows, err := config.AppConfig.Db.Query("insert into job_portal.users (username, password) values ('admin', 'admin')")
//		if err != nil {
//			log.Println("Error in inserting admin", err)
//			return
//		}
//
//		log.Println("User added: ", createdUser)
//
//		w.Write([]byte("Admin added"))
//	}
func main() {
	// init the dotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading the dot env file")
		return
	}

	JWT_SECRET_KEY := os.Getenv("SECRET_KEY")
	if JWT_SECRET_KEY == "" {
		panic("No SECRET_KEY set")
	}

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatal("Error connecting to database", err)
		panic(err)
	}

	// db.TestDB(postgres_db)

	// app-wide state
	config.AppConfig.AddSecretKey(JWT_SECRET_KEY)
	config.AppConfig.Connect(client)

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
	http.ListenAndServe("localhost:5000", r)
}
