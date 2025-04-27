package impl

import (
	"context"
	"errors"
	"mnc/internal/handler"
	auth "mnc/internal/middleware"
	"mnc/internal/model"
	"mnc/internal/repository"
	"mnc/internal/service"
)

func NewUserServiceImpl(userRepository *repository.UserRepository, logService service.LogService) service.UserService {
	return &userServiceImpl{UserRepository: *userRepository, logService: logService}
}

type userServiceImpl struct {
	repository.UserRepository
	logService service.LogService
}

func (s *userServiceImpl) Login(ctx context.Context, request model.UserRequest) ([]model.UserLoginResp, error) {

	// NewLogger().Info("from repo request", request)
	status := "SUCCESS"

	user, found, err := s.UserRepository.FindByUsername(ctx, request) // asumsi: return (user model.UserRequest, found bool)
	if !found {
		msg := "user not found or invalid credentials"
		response := errors.New(msg)
		status = "FAILED"
		errMsg := handler.LogRespErr(msg)
		go s.logService.LogActivity(ctx, user.Id, "LOGIN", request, errMsg, status) // asynchronous logging
		return nil, response
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		msg := "failed to generate token"
		response := errors.New(msg)
		status = "FAILED"
		errMsg := handler.LogRespErr(msg)
		go s.logService.LogActivity(ctx, user.Id, "LOGIN", request, errMsg, status) // asynchronous logging
		return nil, response
	}

	err = s.UserRepository.Login(ctx, user.Id, token)
	if err != nil {
		msg := "failed to Login"
		response := errors.New(msg)
		status = "FAILED"
		errMsg := handler.LogRespErr(msg)
		go s.logService.LogActivity(ctx, user.Id, "LOGIN", request, errMsg, status) // asynchronous logging
		return nil, response
	}

	response := model.UserLoginResp{
		Id:        user.Id,
		Username:  user.Username,
		AuthToken: token,
	}

	go s.logService.LogActivity(ctx, user.Id, "LOGIN", request, response, status) // asynchronous logging

	return []model.UserLoginResp{response}, nil
}

func (s *userServiceImpl) Create(ctx context.Context, request model.UserRequest) (model.UserCreateResp, error) {
	// Validasi / mapping bisa ditaruh di sini
	status := "SUCCESS"

	exists, err := s.UserRepository.FindUserExist(ctx, request.Username)
	if err != nil {
		println("masuk ERROR dari FindUserExist")
		return model.UserCreateResp{}, err
	}
	if exists {
		println("masuk USERNAME SUDAH ADA")
		return model.UserCreateResp{}, errors.New("username already exists")
	}

	hashedPassword, err := handler.HashPassword(request.Password)
	if err != nil {
		status = "FAILED"
		msg := "failed to generate password"
		response := errors.New(msg)
		errMsg := handler.LogRespErr(msg)
		go s.logService.LogActivity(ctx, 0, "CREATE USER", request, errMsg, status) // asynchronous logging
		return model.UserCreateResp{}, response
	}

	user := model.UserRequest{
		Username: request.Username,
		Password: hashedPassword,
		Role:     request.Role,
		IsLogin:  0,
	}
	// Simpan ke database
	savedUser, err := s.UserRepository.Create(ctx, user)
	if err != nil {
		status = "FAILED"
		msg := "failed to create user"
		response := errors.New(msg)
		errMsg := handler.LogRespErr(msg)
		go s.logService.LogActivity(ctx, savedUser.Id, "CREATE USER", request, errMsg, status)
		return model.UserCreateResp{}, response
	}

	// Response
	response := model.UserCreateResp{
		Id:       savedUser.Id,
		Username: savedUser.Username,
	}

	go s.logService.LogActivity(ctx, savedUser.Id, "CREATE USER", request, response, status) // asynchronous logging

	return response, nil
}

func (s *userServiceImpl) Logout(ctx context.Context, request model.UserRequest) (model.UserLogoutResp, error) {
	status := "SUCCESS"
	// ambil ulang data user setelah logout
	user, found, err := s.UserRepository.FindByUsername(ctx, request)
	if !found {
		status = "FAILED"
		msg := "user not found"
		response := errors.New(msg)
		errMsg := handler.LogRespErr(msg)
		go s.logService.LogActivity(ctx, 0, "LOGOUT", request, errMsg, status) // asynchronous logging

		return model.UserLogoutResp{}, response
	}

	err = s.UserRepository.ClearUserSession(ctx, request.Username)
	if err != nil {
		return model.UserLogoutResp{}, err
	}

	response := model.UserLogoutResp{
		Id:       user.Id,
		Username: user.Username,
		IsLogin:  0,
	}

	go s.logService.LogActivity(ctx, user.Id, "LOGOUT", request, response, status) // asynchronous logging

	return response, nil
}
