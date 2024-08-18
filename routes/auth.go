package routes

import (
	"RJD02/job-portal/controllers"
	"github.com/go-chi/chi/v5"
)

func AuthRouter(authRouter chi.Router) {
	authRouter.Get("/", controllers.AuthHome)
	authRouter.Post("/login", controllers.Login)
	authRouter.Post("/signup", controllers.Signup)
}
