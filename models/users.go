package models

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Email    		string `gorm:"unique"`
	Password 		string
	Oauth_provider  string
	Oauth_id		string

	// To OAuth: OneToOne mapping
	OAuth OAuth `gorm:"foreignKey:User_id"`

	// To Profiles: OneToOne mapping
	Profiles Profiles `gorm:"foreignkey:User_id"`

	// To Accounts: OneToOne mapping
	Accounts Accounts `gorm:"foreignkey:User_id"`

	// To Posts: OneToMany mapping
	Posts []Posts `gorm:"foreignKey:User_id"`
}
