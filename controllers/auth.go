package controllers

import (
	"net/http"
)

func AuthHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You've hit api route"))
}
