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

	user, found, err := s.UserRepository.FindByID(ctx, request) // asumsi: return (user model.UserRequest, found bool)
	if !found {
		response := errors.New("user not found or invalid credentials")
		status = "FAILED"
		go s.logService.LogActivity(ctx, user.Id, "LOGIN", request, response, status) // asynchronous logging
		return nil, response
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		response := errors.New("failed to generate token")
		status = "FAILED"
		go s.logService.LogActivity(ctx, user.Id, "LOGIN", request, response, status) // asynchronous logging
		return nil, response
	}

	err = s.UserRepository.Login(ctx, user.Id, token)
	if err != nil {
		response := errors.New("failed to Login")
		status = "FAILED"
		go s.logService.LogActivity(ctx, user.Id, "LOGIN", request, response, status) // asynchronous logging
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
		return model.UserCreateResp{}, err
	}
	if exists {
		return model.UserCreateResp{}, errors.New("username already exists")
	}

	hashedPassword, err := handler.HashPassword(request.Password)
	if err != nil {
		status = "FAILED"
		response := errors.New("failed to generate password")
		go s.logService.LogActivity(ctx, 0, "CREATE USER", request, response, status) // asynchronous logging
		return model.UserCreateResp{}, err
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
		response := errors.New("failed to create user")
		go s.logService.LogActivity(ctx, savedUser.Id, "CREATE USER", request, response, status)
		return model.UserCreateResp{}, err
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
	// Misalnya kita hanya ingin update token menjadi kosong atau ubah status login
	err := s.UserRepository.ClearUserSession(ctx, 0)
	if err != nil {
		return model.UserLogoutResp{}, err
	}
	status := "SUCCESS"
	// Misal ambil ulang data user setelah logout
	user, found, err := s.UserRepository.FindByID(ctx, request)
	if !found {
		status = "FAILED"
		response := errors.New("user not found")
		go s.logService.LogActivity(ctx, user.Id, "LOGOUT", request, response, status) // asynchronous logging

		return model.UserLogoutResp{}, err
	}

	response := model.UserLogoutResp{
		Id:       user.Id,
		Username: user.Username,
		IsLogin:  user.IsLogin,
	}

	go s.logService.LogActivity(ctx, user.Id, "LOGOUT", request, response, status) // asynchronous logging

	return response, nil
}
