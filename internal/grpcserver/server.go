package grpcserver

import (
	"context"
	"log"
	"net"

	"vote/app/service"
	pb "vote/proto/voter"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedVoterServiceServer
}

func (s *server) ValidateInviteCode(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	return service.NewInvitationService().VerifyInviteCode(ctx, req)
}

func StartGRPCServer(grpcPort string) {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterVoterServiceServer(s, &server{})

	reflection.Register(s)

	log.Printf("gRPC server listening on :%s", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}