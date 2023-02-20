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
	fmt.Printf("Sum function is invoked with %v\n", req)
	firstNum := req.GetSum().GetFirstNumber()
	secondNum := req.GetSum().GetSecondNumber()
	result := firstNum + secondNum
	res := &calcpb.SumResponse{
		Result: result,
	}
	return res, nil
}

//func (*server) PrimeNumDecomposition(req *calcpb.PrimeNumberDecomposition, stream calcpb.SumService_PrimeNumberDecompositionServer) error {
//	fmt.Printf("Prime decomposition function is invoked with %v\n", req)
//	number := req.GetNumber()
//	var k int32 = 2
//	for number > 1 {
//		if number%k == 0 {
//			res := &calcpb.PrimeNumberDecompositionResponse{
//				Result: k,
//			}
//			stream.Send(res)
//			number = number / k
//		} else {
//			k++
//		}
//	}
//	return nil
//}

func (*server) PrimeNumberDecomposition(req *calcpb.PrimeNumberDecompositionRequest, stream calcpb.SumService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Received PrimeNumberDecomposition RPC: %v\n", req)

	number := req.GetNumber()
	var divisor int32 = 2

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calcpb.PrimeNumberDecompositionResponse{
				Result: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increased to %v\n", divisor)
		}
	}
	return nil
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
