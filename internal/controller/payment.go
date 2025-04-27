package controller

import (
	"context"
	"mnc/config"
	"mnc/internal/middleware"
	"mnc/internal/model"
	"mnc/internal/repository"
	"mnc/internal/service"

	"github.com/gofiber/fiber/v2"
)

func NewPaymentController(paymentService *service.PaymentService, logger service.LogService, config config.Config, userRepository repository.UserRepository) *PaymentController {
	return &PaymentController{PaymentService: *paymentService, LogService: logger, Config: config, UserRepository: userRepository}
}

type PaymentController struct {
	service.PaymentService
	service.LogService
	config.Config
	repository.UserRepository
}

func (controller PaymentController) PaymentRoute(app *fiber.App) {
	app.Post("/v1/api/transfer", middleware.ValidateAPIKey, middleware.AuthenticateJWT(controller.Config, controller.UserRepository), controller.Transfer)
}

func (controller PaymentController) Transfer(c *fiber.Ctx) error {
	var request model.PaymentRequest
	err := c.BodyParser(&request)
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}
	println("masuk control")
	paymentResponse, err := controller.PaymentService.Create(context.Background(), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    paymentResponse,
	})
}
