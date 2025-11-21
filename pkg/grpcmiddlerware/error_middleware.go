package grpcmiddlerware

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
			debug.PrintStack() // print the stack trace or panic stack trace

			err = status.Errorf(codes.Internal, "Internal server error: %v", r)
		}

	}()

	res, err := handler(ctx, req)

	/*
			Request
		    ↓
			┌─────────────────────────────────────────┐
			│ ErrorMiddleware (start)                 │
			│ - Setup defer recover                   │
			│   ↓                                     │
			│  ┌───────────────────────────────────┐  │
			│  │ HelloWorld Handler                │  │
			│  │ - Validation                      │  │
			│  │ - Process business logic          │  │
			│  │ - Return response/error           │  │
			│  └───────────────────────────────────┘  │
			│   ↓                                     │
			│ - Handle error if exists                │
			│ - Wrap error with codes.Internal        │
			└─────────────────────────────────────────┘
			↓
			Response
	*/

	if err != nil {

		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.Unauthenticated {
				return nil, err
			}
		}
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
		// return nil, err // original error from the handlers
	}
	return res, nil
}
