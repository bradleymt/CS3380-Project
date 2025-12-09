package database

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Publisher struct {
	gorm.Model
	StudioName string `gorm:"unique;not null" json:"studio_name"`
	Country    string `gorm:"not null" json:"country"`
}

type Family struct {
	gorm.Model
}

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username" binding:"required"`
	Email    string `gorm:"unique;not null" json:"email" binding:"required,email"`
	Password string `gorm:"not null" json:"password" binding:"required,min=8"`
	IsAdmin  bool   `gorm:"default:false" json:"is_admin"`

	FamilyID *uint
	Family   *Family `gorm:"foreignKey:FamilyID"`

	PublisherID *uint
	Publisher   *Publisher `gorm:"foreignKey:PublisherID"`
}

func (u *User) CheckHash(password string) error {
	bcryptHash := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(bcryptHash, []byte(password))
}
