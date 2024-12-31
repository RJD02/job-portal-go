package routes

import (
	"RJD02/job-portal/controllers"
	"RJD02/job-portal/middleware"

	"github.com/go-chi/chi/v5"
)

func JobRouter(jobRouter chi.Router) {
	jobRouter.With(middleware.AuthMiddleware).Get("/{id}", controllers.GetJob)
	// jobRouter.Get("/{id}", controllers.GetJob)
	jobRouter.With(middleware.AuthMiddleware).Get("/", controllers.GetJobs)
	// jobRouter.Get("/", controllers.GetJobs)
	jobRouter.With(middleware.AuthMiddleware).Post("/", controllers.AddJob)
	jobRouter.With(middleware.AuthMiddleware).Get("/search", controllers.SearchJobs)
}
