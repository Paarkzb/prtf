package authgrpc

import (
	"context"
	"errors"
	"sso/internal/domain/models"
	authservice "sso/internal/services/auth"
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
	SignIn(ctx context.Context, username string, password string) (tokens models.Tokens, err error)
	SignUp(ctx context.Context, username string, email string, password string) (userID uuid.UUID, err error)
	IsAdmin(ctx context.Context, userId uuid.UUID) (isAdmin bool, err error)
	UserIdentity(ctx context.Context, accessToken string) (auth bool, userID uuid.UUID, err error)
	Refresh(ctx context.Context, userID uuid.UUID, refreshToken string) (tokens models.Tokens, err error)
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

	tokens, err := s.auth.SignIn(ctx, in.GetUsername(), in.GetPassword())
	if err != nil {
		if errors.Is(err, authservice.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid username or password")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.SignInResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil

}

func (s *serverAPI) SignUp(ctx context.Context, in *ssov1.SignUpRequest) (resp *ssov1.SignUpResponse, err error) {

	if in.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	uid, err := s.auth.SignUp(ctx, in.GetUsername(), in.GetEmail(), in.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &ssov1.SignUpResponse{UserID: uid.String()}, nil
}

func (s *serverAPI) IsAdminw(ctx context.Context, in *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if in.UserID == "" {
		return nil, status.Error(codes.InvalidArgument, "userID is required")
	}

	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot parse userID")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "failed to check admin status")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func (s *serverAPI) UserIdentity(ctx context.Context, in *ssov1.UserIdentityRequest) (*ssov1.UserIdentityResponse, error) {
	if in.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "accessToken is required")
	}

	auth, userId, err := s.auth.UserIdentity(ctx, in.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication failed")
	}

	return &ssov1.UserIdentityResponse{Auth: auth, UserID: userId.String()}, nil
}

func (s *serverAPI) Refresh(ctx context.Context, in *ssov1.RefreshRequest) (*ssov1.RefreshResponse, error) {
	if in.UserID == "" {
		return nil, status.Error(codes.InvalidArgument, "userID is required")
	}

	if in.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot parse userID")
	}

	tokens, err := s.auth.Refresh(ctx, userID, in.RefreshToken)
	if err != nil {
		if errors.Is(err, authservice.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid userID or refresh token")
		}

		return nil, status.Error(codes.Internal, "failed to refresh token")
	}

	return &ssov1.RefreshResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}
