package main

import (
	"context"
	"errors"
	"log"
	"net"

	v1pb "github.com/smukherj1/k8s-signer/generated/signer/v1"
	"google.golang.org/grpc"
)

type invocationsServer struct {
	v1pb.UnimplementedInvocationsServer
}

func (s *invocationsServer) CreateInvocation(ctx context.Context, in *v1pb.CreateInvocationRequest) (*v1pb.CreateInvocationResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *invocationsServer) GetInvocation(ctx context.Context, in *v1pb.GetInvocationRequest) (*v1pb.Invocation, error) {
	return nil, errors.New("not implemented")
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	v1pb.RegisterInvocationsServer(s, &invocationsServer{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving: %v", err)
	}
}
