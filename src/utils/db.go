package utils

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" // Used in sql.open()
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
		create table items (name text, price int);
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

	sqlStmt := `
	select name, price from items;
	`
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
