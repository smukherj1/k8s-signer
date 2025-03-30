package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	kafka "github.com/segmentio/kafka-go"
	v1pb "github.com/smukherj1/k8s-signer/generated/signer/v1"
	"github.com/smukherj1/k8s-signer/pkg/log"
	"google.golang.org/grpc"
)

type invocationsServer struct {
	v1pb.UnimplementedInvocationsServer

	kw *kafka.Writer
}

func (s *invocationsServer) CreateInvocation(ctx context.Context, in *v1pb.CreateInvocationRequest) (*v1pb.CreateInvocationResponse, error) {
	log.Infof("CreateInvocation %+v", in)
	s.kw.WriteMessages(ctx, kafka.Message{
		Value: bytes.NewBufferString(
			fmt.Sprintf("%+v", in),
		).Bytes(),
	})
	return &v1pb.CreateInvocationResponse{
		Invocation: &v1pb.Invocation{
			Name:   "ooga booga",
			Params: in.GetParams(),
		},
	}, nil
}

func (s *invocationsServer) GetInvocation(ctx context.Context, in *v1pb.GetInvocationRequest) (*v1pb.Invocation, error) {
	log.Infof("GetInvocation %+v", in)
	return nil, errors.New("not implemented")
}

func (s *invocationsServer) ListInvocations(ctx context.Context, in *v1pb.ListInvocationsRequest) (*v1pb.ListInvocationsResponse, error) {
	log.Infof("ListInvocations %+v", in)
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

func (c *invocationsClient) ListInvocations(ctx context.Context, in *v1pb.ListInvocationsRequest, opts ...grpc.CallOption) (*v1pb.ListInvocationsResponse, error) {
	return c.s.ListInvocations(ctx, in)
}

func runGRPCServer(s *invocationsServer, wg *sync.WaitGroup) {
	defer wg.Done()
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Errorf("Failed to listen for GRPC server: %v", err)
		os.Exit(1)
	}
	srv := grpc.NewServer()
	v1pb.RegisterInvocationsServer(srv, s)
	log.Infof("GRPC server listening at %v.", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		log.Errorf("Error while in GRPC server: %v", err)
		os.Exit(1)
	}
}

func runHTTPGateway(s *invocationsServer, wg *sync.WaitGroup) {
	defer wg.Done()
	ctx := context.Background()
	mux := gwruntime.NewServeMux()
	if err := v1pb.RegisterInvocationsHandlerClient(ctx, mux, &invocationsClient{s: s}); err != nil {
		log.Errorf("Error registering HTTP Gateway: %v", err)
		os.Exit(1)
	}

	addr := ":8080"
	log.Infof("HTTP Gateway listening at %v.", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Errorf("Error serving HTTP Gateway on %v: %v", addr, err)
		os.Exit(1)
	}
}

func main() {
	var wg sync.WaitGroup
	kw := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"redpanda-0.redpanda.redpanda.svc.cluster.local:9093", "redpanda-1.redpanda.redpanda.svc.cluster.local:9093"},
		Topic:    "test-topic",
		Balancer: &kafka.LeastBytes{},
	})
	s := invocationsServer{
		kw: kw,
	}
	wg.Add(2)
	go runGRPCServer(&s, &wg)
	go runHTTPGateway(&s, &wg)
	wg.Wait()
}
