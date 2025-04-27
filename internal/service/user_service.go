package service

import (
	"context"
	"mnc/internal/model"
)

type UserService interface {
	Login(ctx context.Context, request model.UserRequest) ([]model.UserLoginResp, error)
	Create(ctx context.Context, request model.UserRequest) (model.UserCreateResp, error)
	Logout(ctx context.Context, request model.UserRequest) (model.UserLogoutResp, error)
}
