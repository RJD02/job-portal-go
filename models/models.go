package models

import (
	"time"
)

type User struct {
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"last_modified"`
	Email        string    `json:"email"`
}

type Job struct {
	CompanyName string    `json:"company_name"`
	Created     time.Time `json:"created"`
	Img         string    `json:"img"`
	Description string    `json:"description"`
	Role        string    `json:"role"`
	Id          string    `json:"id"`
}

type Response struct {
	Message      string      `json:"message"`
	Error        string      `json:"error"`
	Data         interface{} `json:"data"`
	ResponseCode int         `json:"statuscode"`
}
