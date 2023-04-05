package grpc

import (
	"log"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Open() (*grpc.ClientConn, error) {
	creds, err := credentials.NewClientTLSFromFile("localhost.crt", "")
	if err != nil {
		log.Printf("Failed to load client credentials: %v", err)
	}

	var conn *grpc.ClientConn

	// Set up a connection to the server.

	if err != nil {
		conn, err = grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	} else {
		conn, err = grpc.Dial("localhost:50051", grpc.WithInsecure())
	}

	if err != nil {
		log.Printf("Could not connect: %v", err)
		return nil, err
	}
	defer conn.Close()

	return conn, nil
}
