package routes

import (
	"RJD02/job-portal/controllers"
	"github.com/go-chi/chi/v5"
)

func JobRouter(jobRouter chi.Router) {
	jobRouter.Get("/", controllers.GetJobs)
	jobRouter.Post("/", controllers.AddJob)
}
