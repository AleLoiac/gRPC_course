package main

import (
	"context"
	"fmt"
	"gRPC_course/calculator/calcpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	calcpb.UnimplementedSumServiceServer
}

func (*server) Sum(ctx context.Context, req *calcpb.SumRequest) (*calcpb.SumResponse, error) {
	fmt.Printf("Greet function is invoked with %v\n", req)
	firstNum := req.GetSum().GetFirstNumber()
	secondNum := req.GetSum().GetSecondNumber()
	result := firstNum + secondNum
	res := &calcpb.SumResponse{
		Result: result,
	}
	return res, nil
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calcpb.RegisterSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
