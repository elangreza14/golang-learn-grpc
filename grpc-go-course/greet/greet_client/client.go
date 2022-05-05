package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/elangreza14/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello Im Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could Not Connet To Serve: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// doUnary(c)
	// doServerStream(c)
	// doClientStreaming(c)
	// doBidiStreaming(c)
	doUnaryDeadline(c, 5*time.Second) // sould complete
	doUnaryDeadline(c, 1*time.Second) // sould timeout
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("start unary rpc")
	req := &greetpb.GreetingRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "test",
			LastName:  "kaka",
		},
	}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("err : %v", err)
	}
	log.Printf("response from greet: %v", res)
}

func doServerStream(c greetpb.GreetServiceClient) {
	fmt.Println("start client stream rpc")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "test",
			LastName:  "server streaming",
		},
	}
	streamRes, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("err when calling greetmanytime service : %v", err)
	}
	for {
		res, err := streamRes.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("err when reading greetmanytime service : %v", err)
		}
		fmt.Printf("res: %v\n", res.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "test",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "test 2",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "test 4",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("err when calling longgreet service : %v\n", err)
	}
	for _, req := range requests {
		fmt.Printf("sending request : %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	result, err2 := stream.CloseAndRecv()
	if err2 != nil {
		log.Fatalf("err while receiving longgreet from sever service : %v\n", err)
	}
	fmt.Printf("result : %v\n", result)

}

func doBidiStreaming(c greetpb.GreetServiceClient) {
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error whiel streamin content: %v\n", err)
		return
	}
	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "test sdwewe",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "test 2wewewer",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "test 4retrt",
			},
		},
	}
	waitc := make(chan struct{})

	go func() {
		for _, request := range requests {
			fmt.Printf("Sending Messafge: %v\n ", request)
			stream.Send(request)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	go func() {
		for {

			res, err := stream.Recv()
			if err == io.EOF {
				// log.Fatalf("Error whiel receiving content: %v", err)
				break
			}
			if err != nil {
				log.Fatalf("Error whiel receiving content: %v\n", err)
				break
			}
			fmt.Printf("Result: %v\n", res.GetResult())
		}
		// close(waitc)
		close(waitc)
	}()

	<-waitc
}

func doUnaryDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "test",
			LastName:  "kaka",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		// log.Fatalf("Error While Calling Greet With Deadline")
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit")
			} else {
				log.Fatalf("unexpected error: %v", statusErr)
			}
		} else {
			log.Fatalf("error while calling greet with deadline rpc: %v", err)
		}
		return
	}
	log.Printf("Response from greeting deadline : %v", res.GetResult())

}
