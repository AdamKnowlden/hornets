package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"rpctest/lib/web"

	"github.com/boltdb/bolt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	keys "rpctest/lib/context"
	grpc_dag "rpctest/lib/grpc"
	grpc_merkle "rpctest/lib/grpc/merkle"
)

var Context context.Context

func main() {
	ctx := context.Background()

	wg := new(sync.WaitGroup)

	webFlag := flag.Bool("web", false, "Launch web server: true/false")
	flag.Parse()

	if *webFlag {
		wg.Add(1)
		fmt.Println("Starting with web server enabled")
		go func() {
			err := web.StartServer()

			if err != nil {
				fmt.Println("Fatal error occured in web server")
			}

			wg.Done()
		}()
	}

	ctx = InitDatabase(ctx)
	ctx = InitGrpcServer(ctx)

	Context = ctx

	wg.Wait()
}

func InitDatabase(ctx context.Context) context.Context {
	blockDb, err := bolt.Open("blocks.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	ctx = context.WithValue(ctx, keys.BlockDatabase, blockDb)
	defer blockDb.Close()

	contentDb, err := bolt.Open("content.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	ctx = context.WithValue(ctx, keys.ContentDatabase, contentDb)
	defer contentDb.Close()

	cacheDb, err := bolt.Open("cache.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	ctx = context.WithValue(ctx, keys.CacheDatabase, cacheDb)
	defer cacheDb.Close()

	return ctx
}

func InitGrpcServer(ctx context.Context) context.Context {
	// Create a gRPC server
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile("localhost.crt", "localhost.key")
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_merkle.InterceptContext(ctx)),
	}

	if err == nil {
		opts = append(opts, grpc.Creds(creds))
	} else {
		fmt.Println("Failed to create TLS credentials from file, grpc server will start without TLS and be unsecure")
	}

	s := grpc.NewServer(opts...)
	grpc_dag.RegisterMerkleServiceServer(s, &grpc_merkle.Server{})

	ctx = context.WithValue(ctx, keys.GrpcServer, s)

	// Start the gRPC server
	log.Println("Starting gRPC server...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	defer s.Stop()

	return ctx
}
