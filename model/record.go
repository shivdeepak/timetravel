package model

import (
	"time"

	"gorm.io/gorm"
)

type Record struct {
	gorm.Model

	FirstName  string
	MiddleName string
	LastName   string
	Email      string
	DOB        time.Time
	Phone      string
	Street     string
	City       string
	State      string
	Zip        string
	Country    string
}
