package controllers

import (
	"RJD02/job-portal/config"
	"RJD02/job-portal/db"
	"RJD02/job-portal/models"
	"RJD02/job-portal/utils"
	"context"
	"encoding/json"

	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetJobs(w http.ResponseWriter, r *http.Request) {
	type goodJob struct {
		Data  []db.JobModel `json:"jobs"`
		Total string        `json:"total"`
	}
	var response models.Response
	// get start and end from query params
	startStr := r.URL.Query().Get("start")
	maxResultsStr := r.URL.Query().Get("maxResult")
	start := 0
	maxResults := 10
	if startStr != "" {
		var err error
		start, err = strconv.Atoi(startStr)
		if err != nil {

			response.Message = "start property is not set correctly"
			response.Error = err.Error()
			response.ResponseCode = http.StatusBadRequest
			utils.HandleResponse(w, response)
			return
		}
	}

	if maxResultsStr != "" {
		var err error
		maxResults, err = strconv.Atoi(maxResultsStr)
		if err != nil {
			response.ResponseCode = http.StatusBadRequest
			response.Error = err.Error()
			response.Message = "maxResults property is not set correctly"
			utils.HandleResponse(w, response)
			return
		}
	}

	// get 10 jobs
	jobs, err := config.AppConfig.Db.Job.
		FindMany().
		OrderBy(db.Job.LastModified.Order(db.DESC)).
		Skip(start).
		Take(maxResults).
		Exec(context.Background())

	if err != nil {
		response.Message = "Something went wrong when fetching jobs"
		response.ResponseCode = http.StatusInternalServerError
		response.Error = err.Error()
		utils.HandleResponse(w, response)
		return
	}

	var countResult []struct {
		Count string `json:"count"`
	}
	if err := config.AppConfig.Db.
		Prisma.
		QueryRaw("SELECT COUNT(*) AS count FROM \"Job\"").
		Exec(context.Background(), &countResult); err != nil {
		response.ResponseCode = http.StatusInternalServerError
		response.Error = err.Error()
		response.Message =
			"Something went wrong when getting total count of jobs"
		utils.HandleResponse(w, response)
		return
	}

	if len(jobs) == 0 {
		response.ResponseCode = http.StatusOK
		response.Message = "No jobs found"
		utils.HandleResponse(w, response)
		return
	}

	response.Message = "Successfully fetched the jobs"
	response.ResponseCode = http.StatusOK
	response.Data = goodJob{
		Data:  jobs,
		Total: countResult[0].Count,
	}
	utils.HandleResponse(w, response)

}

// should only be accessed by admin
func AddJob(w http.ResponseWriter, r *http.Request) {
	var job_ db.JobModel
	var response models.Response

	role, ok := r.Context().Value("role").(string)
	if !ok || role != "admin" {
		response.ResponseCode = http.StatusForbidden
		response.Message = "Aww, this route is for admins only"
		utils.HandleResponse(w, response)
		return
	}

	// check if the admin_secret_key is present in the request

	err := json.NewDecoder(r.Body).Decode(&job_)

	if err != nil {
		response.ResponseCode = http.StatusBadRequest
		response.Message = "Error in decoding request body"
		response.Error = err.Error()
		utils.HandleResponse(w, response)
		return
	}

	addedJob, err := config.AppConfig.Db.Job.CreateOne(
		db.Job.CompanyName.Set(job_.CompanyName),
		db.Job.Img.Set(job_.Img),
		db.Job.Description.Set(job_.Description),
		db.Job.Role.Set(job_.Role),
		db.Job.ShortDescription.Set(job_.ShortDescription),
		db.Job.Salary.Set(job_.Salary),
	).Exec(context.Background())

	if err != nil {
		response.ResponseCode = http.StatusInternalServerError
		response.Message = "Failed creating the job"
		response.Error = err.Error()
		utils.HandleResponse(w, response)
		return
	}

	response.ResponseCode = http.StatusOK
	response.Data = addedJob
	response.Message = "Successfully added the job"
	utils.HandleResponse(w, response)
}

func GetJob(w http.ResponseWriter, r *http.Request) {
	jobId := chi.URLParam(r, "id")
	var response models.Response

	if jobId == "" {
		response.ResponseCode = http.StatusNotFound
		response.Message = "job id was not passed"
		utils.HandleResponse(w, response)
		return
	}

	dbJob, err := config.AppConfig.Db.Job.FindUnique(
		db.Job.ID.Equals(jobId),
	).Exec(context.Background())

	if err != nil {
		if ok := db.IsErrNotFound(err); ok {
			response.ResponseCode = http.StatusNotFound
			response.Message = "Job not found"
			response.Error = err.Error()
			utils.HandleResponse(w, response)
			return
		}
		response.Message = "Error getting job"
		response.Error = err.Error()
		response.ResponseCode = http.StatusInternalServerError
		utils.HandleResponse(w, response)
		return
	}

	if dbJob == nil {
		response.Message = "Job not found with this id"
		response.ResponseCode = http.StatusNotFound
		utils.HandleResponse(w, response)
		return
	}

	response.Message = "Found job"
	response.ResponseCode = http.StatusOK
	response.Data = dbJob
	utils.HandleResponse(w, response)
}
