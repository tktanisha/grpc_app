package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/khand/grpc_app/pb"
	"google.golang.org/grpc"
)

type server struct {
	mu sync.Mutex
	pb.UnimplementedUserServiceServer
	users map[string]*pb.User
}

func newServer() *server {
	return &server{
		users: make(map[string]*pb.User),
	}
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if req == nil || req.User == nil {
		return nil, errors.New("invalid request")
	}
	id := uuid.New().String()
	u := &pb.User{Id: id, Name: req.User.Name, Email: req.User.Email}
	s.users[id] = u
	return &pb.CreateUserResponse{User: u}, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[req.Id]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", req.Id)
	}
	return &pb.GetUserResponse{User: u}, nil
}

func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if req == nil || req.User == nil || req.User.Id == "" {
		return nil, errors.New("invalid request")
	}
	u, ok := s.users[req.User.Id]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", req.User.Id)
	}
	u.Name = req.User.Name
	u.Email = req.User.Email
	return &pb.UpdateUserResponse{User: u}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.users[req.Id]
	if !ok {
		return &pb.DeleteUserResponse{Ok: false}, nil
	}
	delete(s.users, req.Id)
	return &pb.DeleteUserResponse{Ok: true}, nil
}

func (s *server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	list := make([]*pb.User, 0, len(s.users))
	for _, u := range s.users {
		list = append(list, u)
	}
	return &pb.ListUsersResponse{Users: list}, nil
}

func main() {
	port := flag.Int("port", 50051, "gRPC server port")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	srv := newServer()
	pb.RegisterUserServiceServer(s, srv)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
