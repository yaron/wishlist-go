package utils

// WishlistItem is a struct that represents an item on a wishlist
type WishlistItem struct {
	ID        int    `json:"id"`
	Name      string `json:"name" binding:"required"`
	Price     int    `json:"price"`
	Claimed   bool   `json:"claimed"`
	Claimable bool   `json:"claimable"`
	URL       string `json:"url"`
	Image     string `json:"image"`
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

// Claim contains all the information needed to claim an item
type Claim struct {
	ID   int    `json:"id" binding:"required"`
	Mail string `json:"mail"`
}
