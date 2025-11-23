// middleware in here because have jwt entity in internal/entity/jwt
package grpcmiddlerware

import (
	"context"

	"google.golang.org/grpc"

	jwtentity "github.com/arthurhzna/Golang_gRPC/internal/entity/jwt"
	"github.com/arthurhzna/Golang_gRPC/internal/utils"
	"github.com/patrickmn/go-cache"
)

type authMiddleware struct {
	cacheService *cache.Cache
}

func NewAuthMiddleware(cacheService *cache.Cache) *authMiddleware {
	return &authMiddleware{
		cacheService: cacheService,
	}
}

var publicApis = map[string]bool{
	"/auth.AuthService/Login":               true,
	"/auth.AuthService/Register":            true,
	"/product.ProductService/DetailProduct": true,
	"/product.ProductService/ListProduct":   true,
}

func (am *authMiddleware) Middleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	if publicApis[info.FullMethod] { // this path get from info.FullMethod
		return handler(ctx, req)
	} // allow login and register without authentication jwt

	// log.Printf("AuthMiddleware: Method=%s", info.FullMethod)

	jwtToken, err := jwtentity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	_, ok := am.cacheService.Get(jwtToken)
	if ok {
		return nil, utils.UnaunthorizedResponse()
	}

	claims, err := jwtentity.GetClaimsFromToken(jwtToken)
	if err != nil {
		// log.Printf("AuthMiddleware: error=%v", err)
		return nil, err
	}

	// log.Printf("AuthMiddleware: claims=%+v", claims)

	ctx = claims.SetToContext(ctx) // store the claims in the context

	res, err := handler(ctx, req)

	// log.Printf("AuthMiddleware: res=%+v", res)

	return res, nil
}
