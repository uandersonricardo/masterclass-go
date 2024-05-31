package internal

import (
	"context"
	"net"

	"github.com/uandersonricardo/masterclass-go/pkg/pb"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	address string
	server  *grpc.Server

	pb.UnimplementedFrameServiceServer
}

func NewGrpcServer(address string) *GrpcServer {
	server := grpc.NewServer()

	return &GrpcServer{
		address: address,
		server:  server,
	}
}

func (s *GrpcServer) Start() error {
	pb.RegisterFrameServiceServer(s.server, s)
	lis, err := net.Listen("tcp", s.address)

	if err != nil {
		return err
	}

	return s.server.Serve(lis)
}

func (s *GrpcServer) GetFrame(ctx context.Context, req *pb.GetFrameRequest) (*pb.Frame, error) {
	return &pb.Frame{
		Id: req.Id,
	}, nil
}
