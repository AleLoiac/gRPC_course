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

	//optional user prompt
	var x, y int32
	_, err = fmt.Scan(&x, &y)
	if err != nil {
		return
	}

	//doUnary(c, x, y)

	doServerStreaming(c, x)
}

func doUnary(c calcpb.SumServiceClient, x int32, y int32) {
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

//func doServerStreaming(c calcpb.SumServiceClient, x int32) {
//	fmt.Println("Starting Server Streaming RPC...")
//	req := &calcpb.PrimeNumberDecompositionRequest{
//		Number: &calcpb.PrimeNumberDecomposition{
//			Number: x,
//		},
//	}
//	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
//	if err != nil {
//		log.Fatalf("Error while calling PrimeDecomposition RPC: %v", err)
//	}
//	for {
//		msg, err := resStream.Recv()
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			log.Fatalf("Error while reading the stream: %v", err)
//		}
//		log.Printf("Response from PrimeNumberDecomposition: %v", msg.GetResult())
//	}
//}

func doServerStreaming(c calcpb.SumServiceClient, x int32) {
	fmt.Println("Starting to do a PrimeDecomposition Server Streaming RPC...")
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
