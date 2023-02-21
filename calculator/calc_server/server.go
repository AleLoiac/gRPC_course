package main

import (
	"context"
	"fmt"
	"gRPC_course/calculator/calcpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct {
	calcpb.UnimplementedCalculatorServiceServer
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

func (*server) PrimeNumberDecomposition(req *calcpb.PrimeNumberDecompositionRequest, stream calcpb.CalculatorService_PrimeNumberDecompositionServer) error {
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

func (*server) ComputeAverage(stream calcpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function is invoked with a streaming request\n")
	var result float64
	var counter int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// EOF is End Of File, it means we have received all the requests that the client sent
			result = result / float64(counter)
			return stream.SendAndClose(&calcpb.ComputeAverageResponse{
				Average: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		number := req.GetNumber()
		result += float64(number)
		counter++
	}
}

func (*server) FindMaximum(stream calcpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function is invoked with a streaming request\n")
	var max int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading the stream: %v", err)
			return err
		}
		number := req.GetNumber()
		if number > max {
			max = number
			sendErr := stream.Send(&calcpb.FindMaximumResponse{
				Maximum: number,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending data to client: %v", sendErr)
				return sendErr
			}
		}
	}
}

func main() {
	fmt.Println("Server started...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calcpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
