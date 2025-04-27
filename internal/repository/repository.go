package repository

import (
	"context"
	"mnc/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, request model.UserRequest) (model.UserCreateResp, error)
	FindByUsername(ctx context.Context, request model.UserRequest) (model.User, bool, error)
	ClearUserSession(ctx context.Context, username string) error
	Login(ctx context.Context, userID int, token string) error
	FindByID(ctx context.Context, userId int) (model.User, bool, error)
	FindUserExist(ctx context.Context, username string) (bool, error)
}

type ActivityLogRepository interface {
	Create(ctx context.Context, log model.ActivityLog) error
}

type PaymentRepository interface {
	Create(ctx context.Context, request model.PaymentRequest) (model.PaymentResp, error)
}
