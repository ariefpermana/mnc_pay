package config

import (
	"mnc/msg"

	"github.com/gofiber/fiber/v2"
)

func NewFiberConfiguration() fiber.Config {
	return fiber.Config{
		ErrorHandler: msg.ErrorHandler,
	}
}
