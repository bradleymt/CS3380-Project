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

type Library struct {
	gorm.Model
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
}

type LibraryItem struct {
	gorm.Model
	LibraryID uint
	Library   Library `gorm:"foreignKey:LibraryID"`
	GameID    uint
	Game      Game `gorm:"foreignKey:GameID"`
}

type Wishlist struct {
	gorm.Model
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
}

type WishlistItem struct {
	gorm.Model
	WishlistID uint
	Wishlist   Wishlist `gorm:"foreignKey:WishlistID"`
	GameID     uint
	Game       Game `gorm:"foreignKey:GameID"`
}

type Cart struct {
	gorm.Model
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
}

type CartItem struct {
	gorm.Model
	CartID uint
	Cart   Cart `gorm:"foreignKey:CartID"`
	GameID uint
	Game   Game `gorm:"foreignKey:GameID"`
}
