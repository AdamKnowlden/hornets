package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	hornet_bolt "rpctest/lib/database/bolt"
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

	//ctx = InitDatabase(ctx)
	//ctx = InitGrpcServer(ctx)

	// Database
	boltDatabase, err := bolt.Open("database.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	boltDb := &hornet_bolt.BoltDatabase{
		Db: boltDatabase,
	}

	boltDb.CreateBucket("blocks")
	boltDb.CreateBucket("content")
	boltDb.CreateBucket("cache")

	ctx = context.WithValue(ctx, keys.BlockDatabase, boltDb)
	defer boltDatabase.Close()

	// Grpc Server
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	customLis := &customListener{Listener: lis}

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
	if err := s.Serve(customLis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	defer s.Stop()

	Context = ctx

	wg.Wait()
}

func InitDatabase(ctx context.Context) context.Context {

	return ctx
}

func InitGrpcServer(ctx context.Context) context.Context {

	return ctx
}

type customListener struct {
	net.Listener
	clientCounter uint64
}

func (cl *customListener) Accept() (net.Conn, error) {
	conn, err := cl.Listener.Accept()
	if err != nil {
		return nil, err
	}

	clientID := atomic.AddUint64(&cl.clientCounter, 1)
	fmt.Printf("New client connected, assigned ID: %d\n", clientID)
	conn = &customConn{Conn: conn, clientID: clientID}
	return conn, nil
}

type customConn struct {
	net.Conn
	clientID uint64
}

func (cc *customConn) ClientID() uint64 {
	return cc.clientID
}
