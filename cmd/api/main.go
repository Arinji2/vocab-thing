package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arinji2/vocab-thing/internal/database"
	"github.com/arinji2/vocab-thing/routes"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db, err := database.SetupDatabase("../../db/app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Database setup complete and ready to use.")

	srv := http.Server{
		Addr:    ":8080",
		Handler: routes.RegisterRoutes(db),
	}

	log.Println("Starting server on port 8080")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
