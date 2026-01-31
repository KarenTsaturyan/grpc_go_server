package auth

import (
	"context"

	ssov1 "github.com/KarenTsaturyan/proto_go/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer // it let's us to run code without implementing whole interface
}

func RegisterServerAPI(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

const (
	emptyValue = 0
)

// Req and Res structures are in proto_go
func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	// TODO: Update with real validation logic
	if req.GetEmail() == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"email is required",
		)
	}

	if req.GetPassword() == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"password is required",
		)
	}

	if req.GetAppId() == emptyValue {
		return nil, status.Error(
			codes.InvalidArgument,
			"app_id is required",
		)
	}

	//  auth Service logic

	return &ssov1.LoginResponse{
		Token: req.GetEmail(),
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Logout(
	ctx context.Context,
	req *ssov1.LogoutRequest,
) (*ssov1.LogoutResponse, error) {
	panic("implement me")
}
