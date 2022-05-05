package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/elangreza14/grpc-go-course/calculator/calculatorpb"
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

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// fmt.Println(c)
	// doUnary(c)
	// doServerStream(c)
	// doClientStreaming(c)
	// doBidiStreaming(c)
	doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("start unary rpc")
	req := &calculatorpb.CalculatoringRequest{
		Calculatoring: &calculatorpb.Calculatoring{
			SumOne: 10,
			SumTwo: 3,
		},
	}
	res, err := c.Calculator(context.Background(), req)

	if err != nil {
		log.Fatalf("err : %v", err)
	}
	log.Printf("response from greet: %v", res)

}

func doServerStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("start client stream rpc")
	req := &calculatorpb.CalculatoringPrimeDecomRequest{
		Primenum: 143,
	}

	streamRes, err := c.CalculatorPrimeDecom(context.Background(), req)

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
		fmt.Printf("res: %v\n", stringTest(res.GetResultprime()))
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {

	requests := []*calculatorpb.AverageRequest{
		{
			Renum: 1,
		}, {
			Renum: 2,
		}, {
			Renum: 3,
		}, {
			Renum: 4,
		},
	}
	stream, err := c.Average(context.Background())
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

func stringTest(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

func doBidiStreaming(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.MaxValue(context.Background())
	if err != nil {
		log.Fatalf("Error whiel streamin content: %v\n", err)
		return
	}
	// requests := []*calculatorpb.MaxValueRequest{
	// 	{
	// 		Inputmax: 1,
	// 	},
	// 	{
	// 		Inputmax: 2,
	// 	},
	// 	{
	// 		Inputmax: 4,
	// 	},
	// 	{
	// 		Inputmax: 3,
	// 	},
	// 	{
	// 		Inputmax: 5,
	// 	},
	// }
	waitc := make(chan struct{})

	go func() {
		numbers := []int32{1, 6, 3, 7, 89, 4}
		for _, number := range numbers {
			fmt.Printf("Sending Messafge: %v\n ", numbers)
			stream.Send(&calculatorpb.MaxValueRequest{Inputmax: number})
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
			fmt.Printf("Result: %v\n", res.GetResultmax())
		}
		// close(waitc)
		close(waitc)
	}()

	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("start unary square rpc")
	doErrorCall(c, int32(9))
	doErrorCall(c, int32(4))
	doErrorCall(c, int32(-1))

}

func doErrorCall(c calculatorpb.CalculatorServiceClient, number int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareeRequest{
		NumberInput: number,
	})

	if err != nil {
		respErr, ok := status.FromError(err)
		// lcog.Fatalf("err : %v", err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We Probably sent negative number")
				return
			}
		} else {
			log.Fatalf("Big Error Calling Square Root : %v\n", err)
			return
		}
	}
	fmt.Printf("response of square root: %v\n", res.GetNumberRoot())
}
