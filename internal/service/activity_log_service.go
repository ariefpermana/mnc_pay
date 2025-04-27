package service

import (
	"context"
)

type LogService interface {
	LogActivity(ctx context.Context, UserId int, action string, req, res interface{}, status string)
}
