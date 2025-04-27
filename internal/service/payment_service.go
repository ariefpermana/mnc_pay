package service

import (
	"context"
	"mnc/internal/model"
)

type PaymentService interface {
	Create(ctx context.Context, request model.PaymentRequest) (model.PaymentResp, error)
}
