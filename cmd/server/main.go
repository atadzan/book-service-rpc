package main

import (
	"context"
	"fmt"
	"github.com/atadzan/book-service-rpc/internal"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
)

func main() {
	log.Println("Starting listening on port 8080")
	port := ":8080"

	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Postgres database
	dbURl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", "admin", "testDB",
		"localhost", "5455", "postgres")
	db, err := pgxpool.Connect(context.Background(), dbURl)
	if err != nil {
		log.Fatalf("Failed to connect to db")
	}
	defer db.Close()

	// Init repository
	var repository internal.BookRepository = internal.NewPostgresRepository(db)
	srv := internal.NewRPCServer(repository)
	if err = srv.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
