package utils

import (
	"RJD02/job-portal/config"
	"context"
	"log"
	"time"
)

func UpdateInactiveLinks() {

	query := `
    update "Job"
    set "isActive" = FALSE
    where "deadline" < NOW() and "isActive" = true
    `

	result, err := config.AppConfig.Db.Prisma.ExecuteRaw(query).Exec(context.Background())
	if err != nil {
		log.Fatal("Faced error when updating inactive jobs", err)
	}

	log.Println("Updated ", result.Count, "jobs to inactive")
}

func RunUpdateInactiveLinksDaily(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)

	UpdateInactiveLinks()

	for {
		select {
		case <-ticker.C:
			UpdateInactiveLinks()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
