package main

import (
	"context"
	"fmt"
	v1 "github.com/atadzan/book-service-rpc/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

const serverAddress = "localhost:8080"

func main() {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ctx := context.Background()
	client := v1.NewBookServiceClient(conn)

	// add book
	bookDTO := &v1.Book{
		Title:       "Golang basic",
		Author:      "Mr X",
		Description: "About Golang syntax",
		Language:    "X language",
		FinishTime:  timestamppb.Now(),
	}
	resCreate, err := client.CreateBook(ctx, &v1.CreateBookRequest{
		Book: bookDTO,
	})
	if err != nil {
		errStatus, _ := status.FromError(err)
		fmt.Println(errStatus.Message())
	}
	log.Printf("Book created with bid: %d\n", resCreate.Bid)

	// Retrieve book by id
	resRetrieve, err := client.RetrieveBook(ctx, &v1.RetrieveBookRequest{
		Bid: resCreate.Bid,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Retrieved book: %s", resRetrieve.String())

	// update book
	bookUpdate := &v1.Book{
		Bid:         resCreate.Bid,
		Title:       "updated 2 Golang basic",
		Author:      "updated 2 Mr X",
		Description: "updated 2 About Golang syntax",
		Language:    "updated X language",
		FinishTime:  timestamppb.Now(),
	}
	_, err = client.UpdateBook(ctx, &v1.UpdateBookRequest{
		Book: bookUpdate,
	})
	if err != nil {
		errStatus, _ := status.FromError(err)
		fmt.Println(errStatus.Message())
	}
	log.Printf("Successfully updated: %d\n", bookUpdate.Bid)

	// delete book
	_, err = client.DeleteBook(ctx, &v1.DeleteBookRequest{Bid: bookUpdate.Bid - 1})
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Successfully updated: %d\n", bookUpdate.Bid-1)

	// List all books
	resListBooks, err := client.ListBook(ctx, &v1.ListBookRequest{
		Limit:  10,
		Offset: 0,
	})
	fmt.Println("Book list: ", resListBooks)
}
