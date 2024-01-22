package grpcapp

import (
	"context"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_SSO/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

type JWTInterceptor struct {
	manager         *jwt.Manager
	accessibleRoles map[string][]string
}

// NewJWTInterceptor creates a new instance of JWTInterceptor with the provided JWT manager and accessibleRoles map.
// The JWTInterceptor is used as a gRPC server interceptor to validate JWT tokens and enforce role-based access control.
func NewJWTInterceptor(
	manager *jwt.Manager,
	accessibleRoles map[string][]string,
) *JWTInterceptor {
	return &JWTInterceptor{manager: manager, accessibleRoles: accessibleRoles}
}

// authorize checks whether the user is authorized to access a specific gRPC method based on JWT token claims and accessible roles.
func (i *JWTInterceptor) authorize(ctx context.Context, method string) error {
	_, ok := i.accessibleRoles[method]
	if !ok {
		// everyone can access
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, grpcerror.ErrNoToken.Error())
	}

	if parts := strings.Fields(values[0]); len(parts) < 2 {
		return status.Error(codes.Unauthenticated, grpcerror.ErrInvalidToken.Error())
	}

	accessToken := strings.Fields(values[0])[1]
	claims, err := i.manager.ParseToken(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, grpcerror.ErrInvalidToken.Error())
	}

	for _, role := range i.accessibleRoles[method] {
		if role == claims["role"] {
			return nil
		}
	}

	return status.Errorf(codes.PermissionDenied, grpcerror.ErrForbidden.Error())
}

// Unary returns a gRPC UnaryServerInterceptor that performs authorization checks before allowing the execution
// of a unary gRPC method.
func (i *JWTInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		err := i.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// Stream returns a gRPC StreamServerInterceptor that performs authorization checks before allowing the execution
// of a streaming gRPC method.
func (i *JWTInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := i.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}
