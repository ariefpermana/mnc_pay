package controller

import (
	"context"
	"mnc/config"
	"mnc/internal/middleware"
	"mnc/internal/model"
	"mnc/internal/service"

	"github.com/gofiber/fiber/v2"
)

func NewUserController(userService *service.UserService, logActivity service.LogService, config config.Config) *UserController {
	return &UserController{UserService: *userService, LogService: logActivity, Config: config}
}

type UserController struct {
	service.UserService
	service.LogService
	config.Config
}

func (controller UserController) UserRoute(app *fiber.App) {
	app.Post("/v1/api/login", middleware.ValidateAPIKey, controller.Login)
	app.Post("/v1/api/create", middleware.ValidateAPIKey, controller.Create)
	app.Post("/v1/api/logout", middleware.ValidateAPIKey, middleware.AuthenticateJWT("ROLE_USER", controller.Config), controller.Logout)
}

func (controller UserController) Login(c *fiber.Ctx) error {
	var request model.UserRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	userResponse, err := controller.UserService.Login(context.Background(), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    userResponse,
	})
}

func (controller UserController) Create(c *fiber.Ctx) error {
	var request model.UserRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	userResponse, err := controller.UserService.Create(context.Background(), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    userResponse,
	})
}

func (controller UserController) Logout(c *fiber.Ctx) error {
	var request model.UserRequest
	err := c.BodyParser(&request)
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	userResponse, err := controller.UserService.Logout(context.Background(), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    userResponse,
	})
}
