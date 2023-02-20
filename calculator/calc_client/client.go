package main

import (
	"context"
	"fmt"
	"gRPC_course/calculator/calcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

func main() {

	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials())) //later on it will need to be secured
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := calcpb.NewCalculatorServiceClient(cc)

	//doUnary(c)

	//doServerStreaming(c)

	doClientStreaming(c)
}

func doUnary(c calcpb.CalculatorServiceClient) {
	var x, y int32
	_, err := fmt.Scan(&x, &y)
	if err != nil {
		return
	}
	fmt.Println("Starting Unary RPC...")
	req := &calcpb.SumRequest{
		Sum: &calcpb.Sum{
			FirstNumber:  x,
			SecondNumber: y,
		},
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.Result)
}

func doServerStreaming(c calcpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a PrimeDecomposition Server Streaming RPC...")
	var x int32
	_, err := fmt.Scan(&x)
	if err != nil {
		return
	}
	req := &calcpb.PrimeNumberDecompositionRequest{
		Number: x,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeDecomposition RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetResult())
	}
}

func doClientStreaming(c calcpb.CalculatorServiceClient) {
	fmt.Println("Starting Client Streaming RPC...")

	requests := []*calcpb.ComputeAverageRequest{
		&calcpb.ComputeAverageRequest{
			Number: 5,
		},
		&calcpb.ComputeAverageRequest{
			Number: 8,
		},
		&calcpb.ComputeAverageRequest{
			Number: 10,
		},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage: %v", err)
	}
	// we iterate over the slice and send each message individually
	for _, req := range requests {
		log.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from ComputeAverage: %v", err)
	}
	fmt.Printf("ComputeAverage response: %v\n", res)
}
