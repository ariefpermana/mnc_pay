package impl

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"mnc/internal/handler"
	"mnc/internal/model"
	"mnc/internal/repository"
	"mnc/internal/service"
	"time"
)

func NewPaymentServiceImpl(
	repo repository.PaymentRepository,
	userRepo repository.UserRepository,
	logService service.LogService,
) service.PaymentService {
	return &payementServiceImpl{
		PaymentRepository: repo,
		userRepo:          userRepo,
		logService:        logService,
	}
}

type payementServiceImpl struct {
	repository.PaymentRepository
	logService service.LogService
	userRepo   repository.UserRepository
}

func (s *payementServiceImpl) Create(ctx context.Context, request model.PaymentRequest) (model.PaymentResp, error) {
	status := "SUCCESS"
	user, found, err := s.userRepo.FindByID(ctx, request.UserId)
	if !found {
		status = "FAILED"
		msg := "user not found"
		response := errors.New(msg)
		errMsg := handler.LogRespErr(msg)
		go s.logService.LogActivity(ctx, request.UserId, "PAYMENT", request, errMsg, status) // asynchronous logging

		return model.PaymentResp{}, response
	}
	fmt.Println("user id : %s", user)

	transactionID := fmt.Sprintf("TRX%s%04d", time.Now().Format("20060102"), rand.Intn(10000))
	request.TrxId = transactionID
	request.Status = "SUCCESS"

	payment, err := s.PaymentRepository.Create(ctx, request)
	if err != nil {
		msg := "failed to add payment"
		response := errors.New(msg)
		status = "FAILED"
		errMsg := handler.LogRespErr(msg)
		go s.logService.LogActivity(ctx, request.UserId, "PAYMENT", request, errMsg, status) // asynchronous logging
		return model.PaymentResp{}, response
	}

	go s.logService.LogActivity(ctx, request.UserId, "PAYMENT", request, payment, status) // asynchronous logging

	return payment, nil
}
