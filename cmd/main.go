package main

import (
	"context"
	"database/sql"
	"github/codingMaster/internal/app/controller"
	"github/codingMaster/internal/database"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
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
	// Establish ngrok listener
	listener, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return err
	}

	log.Println("Ingress established at:", listener.URL())

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
		Address: listener.URL(), // Use ngrok URL
	}

	// Start the application using the ngrok listener
	return http.Serve(listener, app.Moul())
}
