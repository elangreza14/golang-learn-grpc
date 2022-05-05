package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/elangreza14/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetingRequest) (*greetpb.GreetingResponse, error) {
	fmt.Printf("greet function was invoked with : %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "hello" + " " + firstName + " " + lastName
	res := &greetpb.GreetingResponse{
		Result: result,
	}
	return res, nil
	// lastname := req.GetGreeting().LastName()
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("greet many Times function was invoked with : %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	for i := 0; i < 10; i++ {
		result := "hello" + " " + firstName + " " + lastName + " " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with streaming request \n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Failed To Read: %v", err)

		}
		firstName := req.GetGreeting().GetFirstName()
		result += "cek " + firstName + ","
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("LongGreet function was invoked with bidi request \n")
	// result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Failed To Read: %v", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "cek " + firstName + ","
		err2 := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if err2 != nil {
			log.Fatalf("Failed While sending data to client: %v", err)
			return err
		}
	}
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetEWithDeadlineesponse, error) {
	fmt.Printf("greet with deadline function was invoked with : %v\n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			//client cancel the greet
			fmt.Println("client cancel the request")
			return nil, status.Error(codes.DeadlineExceeded, "the client cancel the request")
		}
		time.Sleep(1 * time.Second)
	}

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "hello" + " " + firstName + " " + lastName
	res := &greetpb.GreetEWithDeadlineesponse{
		Result: result,
	}
	return res, nil
	// lastname := req.GetGreeting().LastName()
}

func main() {
	fmt.Println("Starting Server...")

	lis, err := net.Listen("tcp", "localhost:50051")

	if err != nil {
		log.Fatalf("Failed To Listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed To Serve: %v", err)
	}

}
