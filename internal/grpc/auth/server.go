package auth

import (
	"context"
	ssov1 "github.com/zaketn/sso-contracts/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email, password string,
		appId int,
	) (token string, err error)

	RegisterNewUser(
		ctx context.Context,
		email, password string,
	) (userId int64, err error)

	IsAdmin(ctx context.Context, userId int64) (isAdmin bool, err error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

const (
	emptyValue = 0
)

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userId, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userId,
	}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.IsAdmin(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return isAdmin, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "the email value is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "the password value is required")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "the app id is required")
	}

	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "the email value is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "the password value is required")
	}

	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.UserId == emptyValue {
		return status.Error(codes.InvalidArgument, "user is is required")
	}

	return nil
}
