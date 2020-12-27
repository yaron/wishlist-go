package utils

// WishlistItem is a struct that represents an item on a wishlist
type WishlistItem struct {
	Name  string `json:"name" binding:"required"`
	Price int    `json:"price"`
}

// User contains all the properties of a user
type User struct {
	Username string
}

// Login contains all the information needed to authenticate a user
type Login struct {
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}
