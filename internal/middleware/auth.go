package middleware

import (
	"mnc/internal/model"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(os.Getenv("APP_SECRET_KEY"))

func GenerateToken(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString(jwtSecret)
}

// ValidateAPIKey adalah middleware untuk memvalidasi API key
func ValidateAPIKey(c *fiber.Ctx) error {
	// Mengambil nilai header x-api-key
	apiKey := c.Get("x-api-key")
	if apiKey == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "API key is missing",
		})
	}

	// Validasi apakah API key sesuai dengan yang diharapkan
	expectedApiKey := os.Getenv("X_API_KEY")
	if apiKey != expectedApiKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid API key",
		})
	}

	// Jika API key valid, lanjutkan ke handler berikutnya
	return c.Next()
}
