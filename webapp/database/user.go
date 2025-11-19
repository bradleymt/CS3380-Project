package database

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Publisher struct {
	gorm.Model
	StudioName string `gorm:"unique;not null"`
	Country    string `gorm:"not null"`
}

type Family struct {
	gorm.Model
}

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`

	FamilyID *uint
	Family   *Family `gorm:"foreignKey:FamilyID"`

	PublisherID *uint
	Publisher   *Publisher `gorm:"foreignKey:PublisherID"`
}

func (u *User) CheckHash(password string) error {
	bcryptHash := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(bcryptHash, []byte(password))
}
