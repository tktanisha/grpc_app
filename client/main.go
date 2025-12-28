package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/khand/grpc_app/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	serverAddr := flag.String("server", "localhost:50051", "gRPC server address")
	cmd := flag.String("cmd", "list", "command: create|get|update|delete|list")
	id := flag.String("id", "", "user id")
	name := flag.String("name", "", "user name")
	email := flag.String("email", "", "user email")
	flag.Parse()

	dialCtx, dialCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dialCancel()
	conn, err := grpc.DialContext(dialCtx, *serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Wait until connection is READY or dialCtx expires.
	waitCtx, waitCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer waitCancel()
	for conn.GetState() != connectivity.Ready {
		if !conn.WaitForStateChange(waitCtx, conn.GetState()) {
			log.Fatalf("connection state did not become READY: %v", conn.GetState())
		}
	}
	c := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch *cmd {
	case "create":
		if *name == "" || *email == "" {
			fmt.Println("name and email required")
			os.Exit(1)
		}
		resp, err := c.CreateUser(ctx, &pb.CreateUserRequest{User: &pb.User{Name: *name, Email: *email}})
		if err != nil {
			log.Fatalf("CreateUser error: %v", err)
		}
		fmt.Printf("created: %+v\n", resp.User)
	case "get":
		if *id == "" {
			fmt.Println("id required")
			os.Exit(1)
		}
		resp, err := c.GetUser(ctx, &pb.GetUserRequest{Id: *id})
		if err != nil {
			log.Fatalf("GetUser error: %v", err)
		}
		fmt.Printf("user: %+v\n", resp.User)
	case "update":
		if *id == "" {
			fmt.Println("id required")
			os.Exit(1)
		}
		resp, err := c.UpdateUser(ctx, &pb.UpdateUserRequest{User: &pb.User{Id: *id, Name: *name, Email: *email}})
		if err != nil {
			log.Fatalf("UpdateUser error: %v", err)
		}
		fmt.Printf("updated: %+v\n", resp.User)
	case "delete":
		if *id == "" {
			fmt.Println("id required")
			os.Exit(1)
		}
		resp, err := c.DeleteUser(ctx, &pb.DeleteUserRequest{Id: *id})
		if err != nil {
			log.Fatalf("DeleteUser error: %v", err)
		}
		fmt.Printf("deleted ok=%v\n", resp.Ok)
	case "list":
		resp, err := c.ListUsers(ctx, &pb.ListUsersRequest{})
		if err != nil {
			log.Fatalf("ListUsers error: %v", err)
		}
		for _, u := range resp.Users {
			fmt.Printf("- %+v\n", u)
		}
	default:
		fmt.Println("unknown cmd")
	}
}
