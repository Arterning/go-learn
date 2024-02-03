package main

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"testing"
)

var db *sql.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	var err error
	db, err = sql.Open("sqlite", "./test.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		os.Exit(1)
	}

	// 创建表
	createTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        age INTEGER
    )
    `
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println("Error creating table:", err)
		os.Exit(1)
	}
}

func shutdown() {
	db.Close()
}

func TestORM_Insert(t *testing.T) {
	o := NewORM(db)
	user := User{Name: "Alice", Age: 30}

	err := o.Insert("users", user)
	if err != nil {
		t.Errorf("Expected nil error, but got: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Errorf("Expected nil error, but got: %v", err)
	}
	if count == 0 {
		t.Errorf("Expected 1 row, but got: %d", count)
	}
}

func TestORM_Update(t *testing.T) {
	o := NewORM(db)
	user := User{Name: "Alice", Age: 30}
	o.Insert("users", user)

	user.Age = 31
	err := o.Update("users", user, "name = ?", user.Name)
	if err != nil {
		t.Errorf("Expected nil error, but got: %v", err)
	}

	var age int
	err = db.QueryRow("SELECT age FROM users WHERE name = ?", user.Name).Scan(&age)
	if err != nil {
		t.Errorf("Expected nil error, but got: %v", err)
	}
	if age != 31 {
		t.Errorf("Expected age 31, but got: %d", age)
	}
}

func TestORM_FindAll(t *testing.T) {
	o := NewORM(db)
	users := []User{
		{Name: "Alice", Age: 38},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}
	for _, user := range users {
		o.Insert("users", user)
	}

	results, err := o.FindAll("users", "age > ? and name = ?", 32, "Alice")
	if err != nil {
		t.Errorf("Expected nil error, but got: %v", err)
	}
	fmt.Println(results)
	if len(results) == 0 {
		t.Errorf("Expected 2 rows, but got: %d", len(results))
	}
}

func TestORM_Delete(t *testing.T) {
	o := NewORM(db)
	user := User{Name: "Alice", Age: 88}
	o.Insert("users", user)

	err := o.Delete("users", "name = ?", user.Name)
	if err != nil {
		t.Errorf("Expected nil error, but got: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users where name = 'Alice'").Scan(&count)
	if err != nil {
		t.Errorf("Expected nil error, but got: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 rows, but got: %d", count)
	}
}
