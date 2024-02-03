package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

type Database struct {
	Data map[string]string
}

func NewDatabase() *Database {
	return &Database{
		Data: make(map[string]string),
	}
}

func (db *Database) Insert(key, value string) {
	db.Data[key] = value
}

func (db *Database) Select(key string) string {
	return db.Data[key]
}

func (db *Database) Delete(key string) {
	delete(db.Data, key)
}

func (db *Database) Update(key, value string) {
	if _, exists := db.Data[key]; exists {
		db.Data[key] = value
	} else {
		fmt.Println("Key not found")
	}
}

func (db *Database) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(db); err != nil {
		return err
	}

	return nil
}

func LoadDatabaseFromFile(filename string) (*Database, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var db Database
	if err := decoder.Decode(&db); err != nil {
		return nil, err
	}

	return &db, nil
}

func main() {
	var db *Database

	// 尝试从文件加载数据库
	if _, err := os.Stat("database.dat"); err == nil {
		fmt.Println("Loading database from file...")
		db, err = LoadDatabaseFromFile("database.dat")
		if err != nil {
			fmt.Println("Error loading database:", err)
			return
		}
	} else {
		fmt.Println("Creating a new database...")
		db = NewDatabase()
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter SQL statement (INSERT/SELECT/DELETE/UPDATE/EXIT): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch {
		case strings.HasPrefix(input, "INSERT"):
			parts := strings.Split(input, " ")
			if len(parts) != 3 {
				fmt.Println("Invalid INSERT statement. Usage: INSERT key value")
				continue
			}
			key, value := parts[1], parts[2]
			db.Insert(key, value)
			fmt.Println("Record inserted.")

		case strings.HasPrefix(input, "SELECT"):
			parts := strings.Split(input, " ")
			if len(parts) != 2 {
				fmt.Println("Invalid SELECT statement. Usage: SELECT key")
				continue
			}
			key := parts[1]
			value := db.Select(key)
			fmt.Printf("Value: %s\n", value)

		case strings.HasPrefix(input, "DELETE"):
			parts := strings.Split(input, " ")
			if len(parts) != 2 {
				fmt.Println("Invalid DELETE statement. Usage: DELETE key")
				continue
			}
			key := parts[1]
			db.Delete(key)
			fmt.Println("Record deleted.")

		case strings.HasPrefix(input, "UPDATE"):
			parts := strings.Split(input, " ")
			if len(parts) != 3 {
				fmt.Println("Invalid UPDATE statement. Usage: UPDATE key value")
				continue
			}
			key, value := parts[1], parts[2]
			db.Update(key, value)
			fmt.Println("Record updated.")

		case input == "EXIT":
			// 保存数据库到文件
			if err := db.SaveToFile("database.dat"); err != nil {
				fmt.Println("Error saving database:", err)
			}
			fmt.Println("Exiting program.")
			return

		default:
			fmt.Println("Invalid command.")
		}
	}
}
