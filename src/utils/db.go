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
		CREATE TABLE items (
			name text, 
			price int, 
			claimed int DEFAULT 0, 
			claimable int DEFAULT 1,
			url string DEFAULT "",
			image string DEFAULT ""
		);
		CREATE TABLE users (username text, password text);
		`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return fmt.Errorf("Unable to create table %v", err)
		}

		if os.Getenv("WISH_DEBUG") == "1" {
			// test2 - test
			sqlStmt = `
			insert into users (username, password) values ("test2", "$2a$10$iudIPXKg0GlCS7O2K9cMHuD1oot16f8VCrCTMlf.iR7sewzEQdz/.");
			`
			_, err = db.Exec(sqlStmt)
			if err != nil {
				return fmt.Errorf("Unable to create test data %v", err)
			}
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

	sqlStmt := "select rowid, name, price, claimable, claimed, url, image from items;"
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return itemList, fmt.Errorf("Unable to query %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var price int
		var claimable int
		var claimed int
		var url string
		var image string
		err = rows.Scan(&id, &name, &price, &claimable, &claimed, &url, &image)
		if err != nil {
			return itemList, fmt.Errorf("Unable to scan record %v", err)
		}
		itemList = append(itemList, WishlistItem{
			ID:        id,
			Name:      name,
			Price:     price,
			Claimed:   claimed == 1,
			Claimable: claimable == 1,
			URL:       url,
			Image:     image,
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

// EditItem updates an existing item in the wishlist
func EditItem(id int, item WishlistItem) error {
	err := checkDB()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return fmt.Errorf("Unable to open %v", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE items SET name=?, price=? WHERE rowid=?;")
	if err != nil {
		return fmt.Errorf("Unable to update record %v", err)
	}
	_, err = stmt.Exec(item.Name, item.Price, id)
	if err != nil {
		return fmt.Errorf("Unable to update record %v", err)
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

// ClaimItem updates an existing item in the wishlist to be claimed
func ClaimItem(claim Claim) error {
	err := checkDB()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return fmt.Errorf("Unable to open %v", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE items SET claimed=true WHERE rowid=? and claimable=1;")
	if err != nil {
		return fmt.Errorf("Unable to update record %v", err)
	}
	_, err = stmt.Exec(claim.ID)
	if err != nil {
		return fmt.Errorf("Unable to update record %v", err)
	}
	return nil
}
