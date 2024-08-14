package db

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"RJD02/job-portal/models"
	"fmt"

	_ "github.com/lib/pq"
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
	host := os.Getenv("host")
	port, err := strconv.Atoi(os.Getenv("port"))
	if err != nil {
		log.Println("port is not convertable to integer, please chec")
		return nil, err
	}
	user := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}
