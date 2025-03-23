package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1pb "github.com/smukherj1/k8s-signer/generated/signer/v1"
	"google.golang.org/grpc"
)

type invocationsServer struct {
	v1pb.UnimplementedInvocationsServer
}

func (s *invocationsServer) CreateInvocation(ctx context.Context, in *v1pb.CreateInvocationRequest) (*v1pb.CreateInvocationResponse, error) {
	log.Printf("CreateInvocation %+v", in)
	return nil, errors.New("not implemented")
}

func (s *invocationsServer) GetInvocation(ctx context.Context, in *v1pb.GetInvocationRequest) (*v1pb.Invocation, error) {
	log.Printf("GetInvocation %+v", in)
	return nil, errors.New("not implemented")
}

type invocationsClient struct {
	s *invocationsServer
}

func (c *invocationsClient) CreateInvocation(ctx context.Context, in *v1pb.CreateInvocationRequest, opts ...grpc.CallOption) (*v1pb.CreateInvocationResponse, error) {
	return c.s.CreateInvocation(ctx, in)
}

func (c *invocationsClient) GetInvocation(ctx context.Context, in *v1pb.GetInvocationRequest, opts ...grpc.CallOption) (*v1pb.Invocation, error) {
	return c.s.GetInvocation(ctx, in)
}

func runGRPCServer(s *invocationsServer, wg *sync.WaitGroup) {
	defer wg.Done()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen for GRPC server: %v", err)
	}
	srv := grpc.NewServer()
	v1pb.RegisterInvocationsServer(srv, s)
	log.Printf("GRPC server listening at %v.", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Error while in GRPC server: %v", err)
	}
}

func runHTTPGateway(s *invocationsServer, wg *sync.WaitGroup) {
	defer wg.Done()
	ctx := context.Background()
	mux := gwruntime.NewServeMux()
	if err := v1pb.RegisterInvocationsHandlerClient(ctx, mux, &invocationsClient{s: s}); err != nil {
		log.Fatalf("Error registering HTTP Gateway: %v", err)
	}

	addr := ":8090"
	log.Printf("HTTP Gateway listening at %v.", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Error serving HTTP Gateway on %v: %v", addr, err)
	}
}

func main() {
	var wg sync.WaitGroup
	s := invocationsServer{}
	wg.Add(2)
	go runGRPCServer(&s, &wg)
	go runHTTPGateway(&s, &wg)
	wg.Wait()
}
