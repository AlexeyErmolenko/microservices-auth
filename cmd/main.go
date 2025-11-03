package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/AlexeyErmolenko/microservices-auth/pkg/user_v1"
)

const grpcPort = 50051

type userServer struct {
	desc.UnimplementedUserV1Server
}

func (s *userServer) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	return &desc.GetResponse{
		User: &desc.User{
			Id: req.GetId(),
			Info: &desc.UserInfo{
				Name:  "Test",
				Email: "test@mail.ru",
				Role:  desc.Role_ROLE_ADMIN,
			},
			CreatedAt: timestamppb.New(time.Now()),
			UpdatedAt: timestamppb.New(time.Now()),
		},
	}, nil
}

func (s *userServer) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("User info: %v", req.GetInfo())
	return &desc.CreateResponse{Id: int64(rand.Int())}, nil
}

func main() {
	fmt.Println(color.RedString("It's gRPC server"))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &userServer{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
