package impl

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"mnc/internal/model"
	"mnc/internal/repository"
	"mnc/internal/service"
	"time"
)

func NewPaymentServiceImpl(repo repository.PaymentRepository) service.PaymentService {
	return &payementServiceImpl{PaymentRepository: repo}
}

type payementServiceImpl struct {
	repository.PaymentRepository
	logService service.LogService
}

func (s *payementServiceImpl) Create(ctx context.Context, request model.PaymentRequest) ([]model.PaymentResp, error) {
	// NewLogger().Info("from repo request", request)
	status := "SUCCESS"

	tokenVal := ctx.Value("authToken")
	tokenStr, ok := tokenVal.(string)
	if !ok || tokenStr == "" {
		response := errors.New("unauthorized: Token is required")
		status = "FAILED"
		go s.logService.LogActivity(ctx, request.UserId, "PAYMENT", request, response, status) // asynchronous logging
		return nil, response
	}

	transactionID := fmt.Sprintf("TRX%s%04d", time.Now().Format("20060102"), rand.Intn(10000))
	request.TrxId = transactionID
	request.Status = "SUCCESS"

	payment, err := s.PaymentRepository.Create(ctx, request)
	if err != nil {
		response := errors.New("failed to generate token")
		status = "FAILED"
		go s.logService.LogActivity(ctx, request.UserId, "PAYMENT", request, response, status) // asynchronous logging
		return nil, response
	}

	go s.logService.LogActivity(ctx, request.UserId, "PAYMENT", request, payment, status) // asynchronous logging

	return payment, nil
}
