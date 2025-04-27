package middleware

import (
	"context"
	"fmt"
	model "mnc/internal/model"
	"mnc/internal/repository"
	"strings"

	configuration "mnc/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AuthenticateJWT(config configuration.Config, userRepo repository.UserRepository) func(*fiber.Ctx) error {
	jwtSecret := config.Get("JWT_SECRET_KEY")

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
				Code:    401,
				Message: "Unauthorized",
				Data:    "Token tidak ditemukan atau format salah",
			})
		}

		// Ambil token dari header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse dan verifikasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan metode signing adalah HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signing tidak dikenali: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
				Code:    401,
				Message: "Unauthorized",
				Data:    "Token tidak valid atau expired",
			})
		}

		// Ambil klaim dan pastikan username ada
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
				Code:    401,
				Message: "Unauthorized",
				Data:    "Klaim tidak valid",
			})
		}

		username, ok := claims["username"].(string)
		if !ok || username == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
				Code:    401,
				Message: "Unauthorized",
				Data:    "Username tidak valid dalam token",
			})
		}
		cekUser := model.UserRequest{
			Username: username,
		}
		// Validasi username ke database
		user, found, err := userRepo.FindByUsername(context.Background(), cekUser)
		if err != nil || !found {
			return c.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
				Code:    401,
				Message: "Unauthorized",
				Data:    "User tidak ditemukan di database",
			})
		}

		if user.Token != tokenString {
			return c.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
				Code:    401,
				Message: "Unauthorized",
				Data:    "Invalid Token",
			})
		}

		// Simpan user ke context
		c.Locals("user", user)
		return c.Next()
	}
}
