package controller

import (
	"context"
	"mnc/config"
	"mnc/internal/middleware"
	"mnc/internal/model"
	"mnc/internal/service"
	"mnc/msg"

	"github.com/gofiber/fiber/v2"
)

func NewPaymentController(paymentService *service.PaymentService, logger service.LogService, config config.Config) *PaymentController {
	return &PaymentController{PaymentService: *paymentService, LogService: logger, Config: config}
}

type PaymentController struct {
	service.PaymentService
	service.LogService
	config.Config
}

func (controller PaymentController) PaymentRoute(app *fiber.App) {
	app.Post("/v1/api/transfer", middleware.ValidateAPIKey, middleware.AuthenticateJWT("ROLE_USER", controller.Config), controller.Transfer)
}

func (controller PaymentController) Transfer(c *fiber.Ctx) error {
	var request model.PaymentRequest
	err := c.BodyParser(&request)

	// Ambil token dari header Authorization
	token := c.Get("Authorization")

	// Tambahkan token ke context
	ctx := context.WithValue(context.Background(), "authToken", token)

	paymentResponse, err := controller.PaymentService.Create(ctx, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal server error",
		})
	}

	msg.PanicLogging(err)
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    paymentResponse,
	})
}
