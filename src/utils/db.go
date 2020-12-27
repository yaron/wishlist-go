package utils

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" // Used in sql.open()
	"golang.org/x/crypto/bcrypt"
)

const dbfile = "./wishlist.db"

func checkDB() error {
	if _, err := os.Stat(dbfile); os.IsNotExist(err) {
		db, err := sql.Open("sqlite3", dbfile)
		if err != nil {
			return fmt.Errorf("Unable to open %v", err)
		}
		defer db.Close()

		sqlStmt := `
		CREATE TABLE items (name text, price int);
		CREATE TABLE users (username text, password text);
		`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return fmt.Errorf("Unable to create table %v", err)
		}
	}

	return nil
}

// FetchItems returns a list of all items in the wishlist
func FetchItems() ([]WishlistItem, error) {
	var itemList []WishlistItem
	err := checkDB()
	if err != nil {
		return itemList, err
	}

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return itemList, fmt.Errorf("Unable to open %v", err)
	}
	defer db.Close()

	sqlStmt := "select name, price from items;"
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return itemList, fmt.Errorf("Unable to query %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var price int
		err = rows.Scan(&name, &price)
		if err != nil {
			return itemList, fmt.Errorf("Unable to scan record %v", err)
		}
		itemList = append(itemList, WishlistItem{
			Name:  name,
			Price: price,
		})
	}

	err = rows.Err()
	if err != nil {
		return itemList, fmt.Errorf("Unable to fetch rows %v", err)
	}

	return itemList, nil
}

// AddItem creates a new item in the wishlist
func AddItem(item WishlistItem) error {
	err := checkDB()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return fmt.Errorf("Unable to open %v", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO items (name, price) VALUES (?, ?);")
	if err != nil {
		return fmt.Errorf("Unable to create record %v", err)
	}
	_, err = stmt.Exec(item.Name, item.Price)
	if err != nil {
		return fmt.Errorf("Unable to create record %v", err)
	}
	return nil
}

// FetchUser returns a user if the username and password match with a record
func FetchUser(usr, pwd string) (User, error) {
	var user User
	err := checkDB()
	if err != nil {
		return user, err
	}

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return user, fmt.Errorf("Unable to open %v", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT username, password from users WHERE username=? LIMIT 1;")
	if err != nil {
		return user, fmt.Errorf("Unable to query %v", err)
	}

	var name string
	var pass string
	err = stmt.QueryRow(usr).Scan(&name, &pass)
	if err != nil {
		return user, fmt.Errorf("No user found that matches the username or password")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(pwd)); err != nil {
		return user, fmt.Errorf("No user found that matches the password")
	}
	return User{
		Username: name,
	}, nil
}
