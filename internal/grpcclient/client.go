package grpcclient

import (
	// "context"
	"log"
	"sync"

	pb "vote/proto/voter"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	conn   *grpc.ClientConn
	client pb.VoterServiceClient
	once   sync.Once
)

func GetVoterServiceClient() pb.VoterServiceClient {
	once.Do(func() {
		var err error
		// connect to gRPC server (adjust address as needed)
		conn, err = grpc.NewClient(
			"localhost:50051",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		)
		if err != nil {
			log.Fatalf("Failed to dial gRPC: %v", err)
		}

		client = pb.NewVoterServiceClient(conn)
		log.Printf("gRPC client connected to localhost:50051")
	})

	return client
}

// Optional: Close the connection (call during graceful shutdown)
func Close() {
	if conn != nil {
		conn.Close()
	}
}
