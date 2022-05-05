package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/elangreza14/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Calculator(ctx context.Context, req *calculatorpb.CalculatoringRequest) (*calculatorpb.CalculatoringResponse, error) {
	fmt.Printf("greet function was invoked with : %v\n", req)
	sumOneBuff := req.GetCalculatoring().GetSumOne()
	sumTwoBuff := req.GetCalculatoring().GetSumTwo()
	result := sumOneBuff + sumTwoBuff
	res := &calculatorpb.CalculatoringResponse{
		Result: result,
	}
	return res, nil
}

func (*server) CalculatorPrimeDecom(req *calculatorpb.CalculatoringPrimeDecomRequest, stream calculatorpb.CalculatorService_CalculatorPrimeDecomServer) error {
	fmt.Printf("greet many Times function was invoked with : %v\n", req)
	curNum := req.GetPrimenum()
	a := testPrime(curNum)

	for i := 0; i < len(a); i++ {
		res := &calculatorpb.CalculatoringPrimeDecomResponse{
			Resultprime: a[i],
		}
		stream.Send(res)
		// time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	fmt.Printf("LongGreet function was invoked with streaming request \n")
	result := int32(0)
	total := float32(0)
	result2 := []int32{}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Result: total,
			})
		}
		if err != nil {
			log.Fatalf("Failed To Read: %v", err)

		}
		getnum := req.GetRenum()
		result2 = append(result2, int32(getnum))
		result += getnum
		total = float32(result) / float32(len(result2))

	}
}

func (*server) MaxValue(stream calculatorpb.CalculatorService_MaxValueServer) error {
	fmt.Printf("LongGreet function was invoked with bidi request \n")
	maxNum := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Failed To Read: %v", err)
			return err
		}
		base := req.GetInputmax()
		result := base
		if result > maxNum {
			maxNum = result
			err2 := stream.Send(&calculatorpb.MaxValueResponse{
				Resultmax: result,
			})
			if err2 != nil {
				log.Fatalf("Failed While sending data to client: %v", err)
				return err
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareeRequest) (*calculatorpb.SquareResponse, error) {
	fmt.Printf("square function was invoked with : %v\n", req)
	number := req.GetNumberInput()

	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument, fmt.Sprintf("Receive Negative Number: %v\n", number),
		)
	}

	return &calculatorpb.SquareResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Starting Listen Port")
	// testPrint :=

	// fmt.Println(testPrime(120))

	lis, err := net.Listen("tcp", "localhost:50051")

	if err != nil {
		log.Fatalf("Failed To Listen: %v", err)
	}
	fmt.Println("Starting Listen Server Calculator...")
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed To Serve: %v", err)
	}
}

func testPrime(n int32) (cek2 []int32) {
	k := 2
	currTot := int(n)
	for currTot > 1 {
		if currTot%k == 0 {
			currTot = currTot / k
			cek2 = append(cek2, int32(k))
		} else {
			k = k + 1
		}
	}
	return
}
