package main

import (
	"mnc/config"
	"mnc/internal/controller"
	repository "mnc/internal/repository"
	service "mnc/internal/service/impl"

	"mnc/msg"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	//setup configuration
	newConfig := config.New()
	database := config.NewDatabase(newConfig)

	userRepository := repository.InitUserRepo(database)
	paymentRepository := repository.InitPaymentRepo(database)
	logRepository := repository.IniLogRepo(database)

	logService := service.NewLogServiceImp(logRepository)
	userService := service.NewUserServiceImpl(&userRepository, logService)
	paymentService := service.NewPaymentServiceImpl(paymentRepository, userRepository, logService)

	//controller
	userController := controller.NewUserController(&userService, logService, newConfig, userRepository)
	paymentController := controller.NewPaymentController(&paymentService, logService, newConfig, userRepository)

	//setup fiber
	app := fiber.New(config.NewFiberConfiguration())
	app.Use(recover.New())
	app.Use(cors.New())

	//routing
	userController.UserRoute(app)
	paymentController.PaymentRoute(app)

	//swagger
	// app.Get("/swagger/*", swagger.HandlerDefault)

	//start app
	serverPort, _ := os.LookupEnv("APP_PORT")
	err := app.Listen(":" + serverPort)
	msg.PanicLogging(err)
}
