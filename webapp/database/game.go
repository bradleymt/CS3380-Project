package database

import "gorm.io/gorm"

type Game struct {
	gorm.Model
	Name        string `json:"name"`
	PublisherID uint
	Publisher   Publisher `gorm:"foreignKey:PublisherID" json:"-"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	Genre       string    `json:"genre"`
}

type Discount struct {
	gorm.Model
	GameID     uint
	Game       Game `gorm:"foreignKey:GameID"`
	Percentage float64
}

type Purchase struct {
	gorm.Model
	UserID        uint
	User          User `gorm:"foreignKey:UserID"`
	GameID        uint
	Game          Game `gorm:"foreignKey:GameID"`
	Street        string
	City          string
	State         string
	ZipCode       string
	Country       string
	PaymentMethod string
}

type Refund struct {
	gorm.Model
	PurchaseID uint
	Purchase   Purchase `gorm:"foreignKey:PurchaseID"`
	Reason     string
	Approved   bool
}

type LibraryItem struct {
	gorm.Model
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
	GameID uint
	Game   Game `gorm:"foreignKey:GameID"`
}

type WishlistItem struct {
	gorm.Model
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
	GameID uint
	Game   Game `gorm:"foreignKey:GameID"`
}

type CartItem struct {
	gorm.Model
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
	GameID uint
	Game   Game `gorm:"foreignKey:GameID"`
}
