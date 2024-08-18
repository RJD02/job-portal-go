package controllers

import (
	"fmt"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestCheckPasswordHash(t *testing.T) {
	password := "something123"
	salting := 14
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), salting)
	hashedPasswordStr := string(hashedPassword)
	if err != nil {
		t.Errorf("Error in hashing(test), %v", err)
		return
	}

	toTestHashedPassword, err := hashPassword(password)

	if err != nil {
		t.Fatalf("Error in getting hash password %v", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(toTestHashedPassword), []byte(password))

	if err != nil {
		fmt.Println("hashedPasswordStr = ", hashedPasswordStr)
		fmt.Println("toTestHashedPassword = ", toTestHashedPassword)
		t.Fatal("Hashes don't match")
		return
	}
}
