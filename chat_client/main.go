package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	pb "github.com/proudlily/test_circleci/proto"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		// Contact the server and print out its response.
		name := scanner.Text()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.SayHello(ctx, &pb.HelloRequest{Message: name})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("%s", r.GetMessage())
	}

}
