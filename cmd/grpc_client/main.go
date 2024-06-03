package main

import (
	"context"
	desc "github.com/Alzoww/go-grpc-gateway-example/pkg/note_v1"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	grpcAddr = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := desc.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.Get(ctx, &desc.GetRequest{Id: 12})
	if err != nil {
		log.Fatalf("could not get note: %v", err)
	}

	log.Printf(color.GreenString("Note info:\n"), color.BlueString("%+v\n", resp.GetNote()))
}
