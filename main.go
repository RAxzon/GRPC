package main

import (
	pb "./proto/consignment"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

const (
	port = ":50051"
)

type repository interface {
	Create(consignment *pb.Consignment) (*pb.Consignment, error)
}

type Repository struct {
	mu sync.RWMutex
	consignments []*pb.Consignment
}

type service struct {
	repo repository
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	repo.mu.Unlock()

	return consignment, nil
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {

	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return &pb.Response{Created: true, Consignment:consignment}, nil
}

func main() {

	repo := &Repository{}

	// Set up GRPC server---------
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	//---------------------------

	// Register service with GRPC server.

	pb.RegisterShippingServiceServer(s, &service{repo})

	// Register reflection service on GRPC server

	log.Printf("Running on port: %v", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}