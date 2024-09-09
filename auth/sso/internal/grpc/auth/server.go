package auth

import (
	"context"
	"errors"
	"sso/internal/services/auth"
	"sso/internal/storage"
	ssov1 "sso/protos/gen/go/sso"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	SignIn(ctx context.Context, username string, password string) (token string, err error)
	SignUp(ctx context.Context, username string, password string) (userID uuid.UUID, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) SignIn(ctx context.Context, in *ssov1.SignInRequest) (resp *ssov1.SignInResponse, err error) {

	if in.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := s.auth.SignIn(ctx, in.GetUsername(), in.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.SignInResponse{Token: token}, nil

}

func (s *serverAPI) SignUp(ctx context.Context, in *ssov1.SignUpRequest) (resp *ssov1.SignUpResponse, err error) {

	if in.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	uid, err := s.auth.SignUp(ctx, in.GetUsername(), in.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &ssov1.SignUpResponse{UserId: uid.String()}, nil
}
