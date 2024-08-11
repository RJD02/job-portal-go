package db

import (
	"database/sql"

	"RJD02/job-portal/models"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func TestDB(db *sql.DB) {
	db.Ping()
	rows, err := db.Query("select * from job_portal.users")
	if err != nil {
		fmt.Println("Error in db query", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.Username, &u.Password, &u.Created, &u.LastModified)
		if err != nil {
			fmt.Println("Error in scanning", err)
			return
		}
		fmt.Println(u.Username, u.Password, u.Created, u.LastModified)
	}
	fmt.Println("Tested db")
}

func Connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}
