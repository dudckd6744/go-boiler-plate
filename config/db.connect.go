package config

import (
	"database/sql"
	"fmt"
	"time"
)

func ConnectionDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:12341234@tcp(localhost:3306)/go_test")

	if err != nil {
		fmt.Printf("Error %s when opening DB\n", err)
	}
	db.SetConnMaxIdleTime(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute * 3)

	_, err = createTables(db)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) (sql.Result, error) {

	query := `CREATE TABLE IF NOT EXISTS User(id int primary key auto_increment, name text,  email varchar(20) ,
		age int, created_at datetime default CURRENT_TIMESTAMP, updated_at datetime default CURRENT_TIMESTAMP)`

	res, err := db.Exec(query)

	if err != nil {
		return nil, err

	}

	return res, nil

}
