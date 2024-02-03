package main

type User struct {
	ID   int    `db:"id" auto:"true"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}
