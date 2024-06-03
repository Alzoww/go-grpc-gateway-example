package main

import (
	"context"
	desc "github.com/Alzoww/go-grpc-gateway-example/pkg/note_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"net/http"
	"sync"
)

const (
	grpcAddress = "localhost:50051"
	httpAddress = "localhost:8080"
)

type Server struct {
	desc.UnimplementedNoteV1Server
}

func (s *Server) Get(ctx context.Context, in *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("Note id: %d", in.GetId())

	return &desc.GetResponse{
		Note: &desc.Note{
			Id: in.GetId(),
			Info: &desc.NoteInfo{
				Title:    gofakeit.BeerName(),
				Context:  gofakeit.IPv4Address(),
				Author:   gofakeit.Name(),
				IsPublic: true,
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func main() {
	ctx := context.Background()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		if err := startGrpcServer(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()

		if err := startHttpServer(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}

func startGrpcServer() error {
	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(grpcServer)

	desc.RegisterNoteV1Server(grpcServer, &Server{})

	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		return err
	}

	log.Printf("grpc server listening on %s", grpcAddress)

	return grpcServer.Serve(listener)
}

func startHttpServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := desc.RegisterNoteV1HandlerFromEndpoint(ctx, mux, grpcAddress, opts); err != nil {
		return err
	}

	log.Printf("http server listening at %v\n", httpAddress)
	return http.ListenAndServe(httpAddress, mux)
}
