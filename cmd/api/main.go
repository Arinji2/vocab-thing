package main

import (
	"fmt"
	"log"

	"github.com/arinji2/vocab-thing/db"
)

func main() {
	fmt.Println("Hello, world!")
	db, err := db.SetupDatabase("../../db/app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Database setup complete and ready to use.")
}
