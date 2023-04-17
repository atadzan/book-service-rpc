package internal

import (
	"context"
	v1 "github.com/atadzan/book-service-rpc/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

type grpcServer struct {
	BookRepository BookRepository
	v1.UnimplementedBookServiceServer
}

func NewRPCServer(repository BookRepository) *grpc.Server {
	srv := grpcServer{
		BookRepository: repository,
	}
	gsrv := grpc.NewServer()

	v1.RegisterBookServiceServer(gsrv, &srv)
	return gsrv
}

func (s *grpcServer) CreateBook(ctx context.Context, req *v1.CreateBookRequest) (*v1.CreateBookResp, error) {
	book := &Book{
		Bid:         0,
		Title:       req.Book.GetTitle(),
		Author:      req.Book.Author,
		Description: req.Book.GetDescription(),
		Language:    req.Book.GetLanguage(),
		FinishTime:  req.Book.GetFinishTime().AsTime(),
	}
	bid, err := s.BookRepository.CreateBook(ctx, book)
	if err != nil {
		log.Println(err.Error())
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &v1.CreateBookResp{Bid: int64(bid)}, nil
}

func (s *grpcServer) RetrieveBook(ctx context.Context, req *v1.RetrieveBookRequest) (*v1.RetrieveBookResp, error) {
	bookId := BookId(req.GetBid())
	book, err := s.BookRepository.RetrieveBook(ctx, bookId)
	if err != nil {
		log.Fatal(err.Error())
	}
	res := &v1.RetrieveBookResp{
		Book: &v1.Book{
			Bid:         int64(book.Bid),
			Title:       book.Title,
			Author:      book.Author,
			Description: book.Description,
			Language:    book.Language,
			FinishTime:  timestamppb.New(book.FinishTime),
		},
	}
	return res, nil
}

func (s *grpcServer) UpdateBook(ctx context.Context, req *v1.UpdateBookRequest) (*v1.UpdateBookResponse, error) {
	book := &Book{
		Bid:         BookId(req.Book.GetBid()),
		Title:       req.Book.GetTitle(),
		Author:      req.Book.GetAuthor(),
		Description: req.Book.GetDescription(),
		Language:    req.Book.GetLanguage(),
	}
	err := s.BookRepository.UpdateBook(ctx, book)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &v1.UpdateBookResponse{}, nil
}

func (s *grpcServer) DeleteBook(ctx context.Context, req *v1.DeleteBookRequest) (*v1.DeleteBookResponse, error) {
	err := s.BookRepository.DeleteBook(ctx, BookId(req.Bid))
	if err != nil {
		log.Fatal(err.Error())
	}
	return &v1.DeleteBookResponse{}, nil
}

func (s *grpcServer) ListBook(ctx context.Context, req *v1.ListBookRequest) (*v1.ListBookResponse, error) {
	books, err := s.BookRepository.ListBook(ctx, req.Offset, req.Limit)
	if err != nil {
		log.Fatal(err.Error())
	}
	res := &v1.ListBookResponse{}
	data := []*v1.Book{}
	for _, book := range books {
		b := &v1.Book{
			Bid:         int64(book.Bid),
			Title:       book.Title,
			Description: book.Description,
			Author:      book.Author,
			Language:    book.Language,
			FinishTime:  timestamppb.New(book.FinishTime),
		}
		data = append(data, b)
	}
	res.Books = data
	return res, nil
}
