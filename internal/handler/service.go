package handler

import (
	"context"
	"fmt"

	"github.com/arthurhzna/Golang_gRPC/internal/utils"
	"github.com/arthurhzna/Golang_gRPC/pb/service"
)

// type IServiceHandler interface {
// service.HelloWorldServiceServer // interface from the service package
/* if im using this --> service.HelloWorldServiceServer is the interface from the service package

// Interface HelloWorldServiceServer (dari generated code) punya method:
type HelloWorldServiceServer interface {
	HelloWorld(ctx context.Context, req *HelloWorldRequest) (*HelloWorldResponse, error)
	mustEmbedUnimplementedHelloWorldServiceServer()
}

// Maka IServiceHandler otomatis punya method yang sama:
type IServiceHandler interface {
	HelloWorld(ctx context.Context, req *HelloWorldRequest) (*HelloWorldResponse, error)
	mustEmbedUnimplementedHelloWorldServiceServer()
}
*/
// 	HelloWorld(ctx context.Context, req *service.HelloWorldRequest) (*service.HelloWorldResponse, error)
// }

type serviceHandler struct {
	service.UnimplementedHelloWorldServiceServer // parent struct from the service package
}

/*
		// Java style
		class UnimplementedHelloWorldServiceServer {
			public HelloWorldResponse HelloWorld(...) {
				throw new Error("not implemented");
			}

			public void mustEmbedUnimplemented...() {
				// ...
			}
		}

		class ServiceHandler extends UnimplementedHelloWorldServiceServer {  // ← INHERIT
			@Override
			public HelloWorldResponse HelloWorld(...) {  // ← OVERRIDE
				// Your implementation
				return new HelloWorldResponse(...);
			}

			// mustEmbedUnimplemented...() tidak di-override, pakai dari parent
		}

	┌────────────────────────────────────────────┐
	│ UnimplementedHelloWorldServiceServer       │  ← "Parent" / Base
	├────────────────────────────────────────────┤
	│ + HelloWorld() → error "not implemented"   │
	│ + mustEmbedUnimplemented...()              │
	└────────────────────────────────────────────┘
						↑
						│ EMBEDDED (seperti inheritance)
						│
	┌────────────────────────────────────────────┐
	│ serviceHandler                             │  ← "Child" / Derived
	├────────────────────────────────────────────┤
	│ • UnimplementedHelloWorldServiceServer     │  ← embedded field
	│                                            │
	│ + HelloWorld() → YOUR IMPLEMENTATION ✅    │  ← OVERRIDE
	│ + mustEmbedUnimplemented...() (inherited)  │  ← TIDAK di-override
	└────────────────────────────────────────────┘
*/

func NewServiceHandler() *serviceHandler {
	return &serviceHandler{}
}

func (sh *serviceHandler) HelloWorld(ctx context.Context, req *service.HelloWorldRequest) (*service.HelloWorldResponse, error) { //override the method from the parent struct

	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &service.HelloWorldResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	return &service.HelloWorldResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
		Base:    utils.SuccessResponse("Success"),
	}, nil
}

// func NewServiceHandler() IServiceHandler {
// 	return &serviceHandler{}
// }
