package main

import (
	"context"
	"fmt"
	"gRPC_course/calculator/calcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

func main() {

	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials())) //later on it will need to be secured
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := calcpb.NewSumServiceClient(cc)

	//doUnary(c)

	doServerStreaming(c)
}

func doUnary(c calcpb.SumServiceClient) {
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

func doServerStreaming(c calcpb.SumServiceClient) {
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
