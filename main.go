package main

import (
	"context"
	"database/sql"
	"github/codingMaster/internal/app/controller"
	"github/codingMaster/internal/database"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var (
	Config = ":3000"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	// Database connection
	db, err := sql.Open("postgres", "postgresql://coding_master_2025_user:wbc6Co8PuuCSDSqzO8fEvMXMjEJ9c9Mx@dpg-cvt1hba4d50c73dbu5og-a.oregon-postgres.render.com/coding_master_2025")
	if err != nil {
		log.Printf("err type: %T\n", err)
		log.Fatal("There is an err in the connection")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Printf("err type: %T\n", err)
		log.Printf("err is: %v\n", err)
		log.Fatal("There is an err in Ping")
	}
	log.Println("Database connection successful")

	st_ore := database.NewStore(db)
	log.Println("Store has been opened")

	app := &controller.Application{
		Store:   st_ore,
		Address: "http://localhost" + Config, // Use local address
	}

	// Start the application on the local port
	log.Println("Server is running on", Config)
	return http.ListenAndServe(Config, app.Moul())
}
