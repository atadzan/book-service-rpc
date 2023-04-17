package internal

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

const booksTable = "books"

type PostgresRepository struct {
	dbPool *pgxpool.Pool
}

func NewPostgresRepository(dbPool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		dbPool: dbPool,
	}
}
func (r *PostgresRepository) CreateBook(ctx context.Context, book *Book) (BookId, error) {
	query := fmt.Sprintf(`INSERT INTO %s (title, author, description, language, finish_time) VALUES($1, $2, $3, $4, $5) RETURNING id`, booksTable)
	var id BookId
	err := r.dbPool.QueryRow(ctx, query, book.Title, book.Author, book.Description, book.Language, book.FinishTime).Scan(&id)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return id, nil
}

func (r *PostgresRepository) RetrieveBook(ctx context.Context, bookId BookId) (*Book, error) {
	var book Book
	query := fmt.Sprintf(`SELECT id, title, description, author, language, finish_time FROM %s WHERE id = $1`, booksTable)
	if err := r.dbPool.QueryRow(ctx, query, bookId).Scan(&book.Bid, &book.Title, &book.Description, &book.Author, &book.Language, &book.FinishTime); err != nil {
		log.Fatalf(err.Error())
	}
	return &book, nil
}

func (r *PostgresRepository) UpdateBook(ctx context.Context, book *Book) error {
	query := fmt.Sprintf(`UPDATE %s SET title = $1, author = $2, description = $3, language = $4, finish_time = now() WHERE id = $5`, booksTable)
	_, err := r.dbPool.Exec(ctx, query, book.Title, book.Author, book.Description, book.Language, book.Bid)
	if err != nil {
		log.Fatal(err.Error())
	}
	return nil
}

func (r *PostgresRepository) DeleteBook(ctx context.Context, bookId BookId) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, booksTable)
	_, err := r.dbPool.Exec(ctx, query, bookId)
	if err != nil {
		log.Fatal(err.Error())
	}
	return nil
}

func (r *PostgresRepository) ListBook(ctx context.Context, offset, limit int64) ([]*Book, error) {
	query := fmt.Sprintf(`SELECT id, title, description, author, language, finish_time FROM %s LIMIT $1 OFFSET $2`, booksTable)
	rows, err := r.dbPool.Query(ctx, query, limit, offset)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()
	var books []*Book
	for rows.Next() {
		var book Book
		if err = rows.Scan(&book.Bid, &book.Title, &book.Description, &book.Author, &book.Language, &book.FinishTime); err != nil {
			log.Fatal(err.Error())
		}
		books = append(books, &book)
	}
	fmt.Println("repo:", books)
	return books, nil
}
