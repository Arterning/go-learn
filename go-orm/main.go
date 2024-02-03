package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
)

func main() {

	db, err := sql.Open("sqlite3", "./test.db")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = createTable(err, db)

	// 创建ORM实例
	o := NewORM(db)

	// 插入数据
	user := User{Name: "Alice", Age: 30}
	err = o.Insert("users", user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data inserted successfully.")
}

func createTable(err error, db *sql.DB) error {
	// 创建表
	createTableSql := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY,
        name TEXT,
        age INTEGER
    )
    `
	_, err = db.Exec(createTableSql)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
