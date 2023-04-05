package merkle

import (
	"context"
	keys "rpctest/lib/context"
	grpc_dag "rpctest/lib/grpc"

	"google.golang.org/grpc"
)

type Server struct {
	grpc_dag.MerkleServiceServer
}

/*
func (s *server) SendMerkleDag(ctx context.Context, in *grpc_dag.MerkleRoot) (*grpc_dag.MerkleDag, error) {
	// Convert the gRPC MerkleDag to your MerkleDag struct
	dag := merkle_dag.FromGRPCMerkleDag(in)

	// Do something with the MerkleDag (e.g., print it to the console)
	fmt.Printf("Received MerkleDag: %+v\n", dag)

	// Convert the MerkleDag struct back to a gRPC MerkleDag
	out := merkle_dag.ToGRPCMerkleDag(dag)

	return out, nil
}
*/

func CreateService(conn *grpc.ClientConn) (*grpc_dag.MerkleServiceClient, error) {
	// Create a new MyServiceClient using the connection.
	client := grpc_dag.NewMerkleServiceClient(conn)

	return &client, nil
}

func (s *Server) SendMerkleRoot(ctx context.Context, in *grpc_dag.MerkleRoot) (*grpc_dag.Response, error) {

	response := &grpc_dag.Response{
		Message: "Root recieved",
	}

	return response, nil
}

func (s *Server) SendMerkleNode(ctx context.Context, in *grpc_dag.MerkleNode) (*grpc_dag.Response, error) {

	response := &grpc_dag.Response{
		Message: "Node recieved",
	}

	return response, nil
}

func (s *Server) NotifyCompletion(ctx context.Context, in *grpc_dag.MerkleRoot) (*grpc_dag.Reciept, error) {

	reciept := &grpc_dag.Reciept{
		Root: in,
	}

	return reciept, nil
}

func (s *Server) SendSignedReciept(ctx context.Context, in *grpc_dag.SignedReciept) (*grpc_dag.Response, error) {

	response := &grpc_dag.Response{
		Message: "Signed reciept recieved",
	}

	return response, nil
}

func InterceptContext(ctx context.Context) grpc.UnaryServerInterceptor {
	return func(newCtx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		blockDb := ctx.Value(keys.BlockDatabase)
		contentDb := ctx.Value(keys.ContentDatabase)
		cacheDb := ctx.Value(keys.CacheDatabase)
		grpcServer := ctx.Value(keys.GrpcServer)

		newCtx = context.WithValue(newCtx, keys.BlockDatabase, blockDb)
		newCtx = context.WithValue(newCtx, keys.ContentDatabase, contentDb)
		newCtx = context.WithValue(newCtx, keys.CacheDatabase, cacheDb)
		newCtx = context.WithValue(newCtx, keys.GrpcServer, grpcServer)

		// Call the next handler with the updated context
		return handler(newCtx, req)
	}
}
