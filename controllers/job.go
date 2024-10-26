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
)

func GetJobs(w http.ResponseWriter, r *http.Request) {
	var response models.Response
	// get start and end from query params
	startStr := r.URL.Query().Get("start")
	start := 0
	if startStr != "" {
		var err error
		start, err = strconv.Atoi(startStr)
		if err != nil {

			http.Error(w, "Invalid start", http.StatusBadRequest)
			return
		}
	}

	// get 10 jobs
	jobs, err := config.AppConfig.Db.Job.
		FindMany().
		OrderBy(db.Job.LastModified.Order(db.DESC)).
		Skip(start).
		Take(10).Exec(context.Background())

	if err != nil {
		response.Message = "Something went wrong when fetching jobs"
		response.ResponseCode = http.StatusInternalServerError
		response.Error = err.Error()
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
	response.Data = jobs
	utils.HandleResponse(w, response)

}

func AddJob(w http.ResponseWriter, r *http.Request) {
	var job models.Job
	var response models.Response

	err := json.NewDecoder(r.Body).Decode(&job)

	if err != nil {
		response.ResponseCode = http.StatusBadRequest
		response.Message = "Error in decoding request body"
		response.Error = err.Error()
		utils.HandleResponse(w, response)
		return
	}

	addedJob, err := config.AppConfig.Db.Job.CreateOne(
		db.Job.CompanyName.Set(job.CompanyName),
		db.Job.Img.Set(job.Img),
		db.Job.Description.Set(job.Description),
		db.Job.Role.Set(job.Role),
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
	utils.HandleResponse(w, response)
}

func GetJob(w http.ResponseWriter, r *http.Request) {
	jobId := r.URL.Query().Get("jobid")
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

//
// func GetJob(w http.ResponseWriter, r *http.Request) {
// 	var j models.Job
// 	id := chi.URLParam(r, "id")
// 	err := config.AppConfig.Db.QueryRow("select id, company_name, created, img, description, role from job_portal.jobs where id = $1", id).Scan(&j.Id, &j.CompanyName, &j.Created, &j.Img, &j.Description, &j.Role)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			http.Error(w, "Job not found", http.StatusNotFound)
// 		} else {
// 			http.Error(w, "Error in getting job", http.StatusInternalServerError)
// 			log.Println("Error in getting job", err)
// 		}
// 		return
// 	}
//
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(j)
//
// }
//
// func AddJob(w http.ResponseWriter, r *http.Request) {
// 	// get the job from request body
// 	var j models.Job
// 	err := json.NewDecoder(r.Body).Decode(&j)
// 	if err != nil {
// 		http.Error(w, "Error in decoding", http.StatusBadRequest)
// 		return
// 	}
//
// 	j.Created = time.Now()
//
// 	j.Id = uuid.New().String()
//
// 	// insert the job
// 	_, err = config.AppConfig.Db.Exec("insert into job_portal.jobs (company_name, img, description, role, created, id) values ($1, $2, $3, $4, current_timestamp, $5)", j.CompanyName, j.Img, j.Description, j.Role, j.Id)
// 	if err != nil {
// 		http.Error(w, "Error in inserting job", http.StatusInternalServerError)
// 		log.Println("Error in inserting job", err)
// 		return
// 	}
//
// 	w.Write([]byte("Job added"))
// }
