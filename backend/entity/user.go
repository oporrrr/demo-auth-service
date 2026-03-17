package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	AccountID     string  `gorm:"uniqueIndex;not null"`
	Role          string  `gorm:"default:'USER'"`
	AccountStatus string  `gorm:"default:'ACTIVE'"`
	FirstName     string
	LastName      string
	Email         string  `gorm:"index"`
	PhoneNumber   string  `gorm:"index"`
	CountryCode   string
	PrefixName    string
	Gender        string
	DateOfBirth   string
	CisNumber     *string
}
