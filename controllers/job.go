package controllers

import (
	"RJD02/job-portal/config"
	"RJD02/job-portal/models"
	"encoding/json"
	"log"
	"time"

	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func GetJobs(w http.ResponseWriter, r *http.Request) {
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
	rows, err := config.AppConfig.Db.Query("select company_name, created, img, description, role, id from job_portal.jobs order by created desc limit 10 offset $1", start)
	if err != nil {
		http.Error(w, "Error in getting jobs", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	jobs := []models.Job{}

	for rows.Next() {
		var j models.Job
		err := rows.Scan(&j.CompanyName, &j.Created, &j.Img, &j.Description, &j.Role, &j.Id)
		if err != nil {
			log.Println("Error in scanning", err)
			http.Error(w, "Error in scanning", http.StatusInternalServerError)
			return
		}
		jobs = append(jobs, j)
	}

	if len(jobs) == 0 {
		http.Error(w, "No jobs found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// return the jobs
	json.NewEncoder(w).Encode(jobs)
}

func AddJob(w http.ResponseWriter, r *http.Request) {
	// get the job from request body
	var j models.Job
	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		http.Error(w, "Error in decoding", http.StatusBadRequest)
		return
	}

	j.Created = time.Now()

	j.Id = uuid.New().String()

	// insert the job
	_, err = config.AppConfig.Db.Exec("insert into job_portal.jobs (company_name, img, description, role, created, id) values ($1, $2, $3, $4, current_timestamp, $5)", j.CompanyName, j.Img, j.Description, j.Role, j.Id)
	if err != nil {
		http.Error(w, "Error in inserting job", http.StatusInternalServerError)
		log.Println("Error in inserting job", err)
		return
	}

	w.Write([]byte("Job added"))
}
