package handlers

import (
	"backend-pedika-fiber/database"
	"backend-pedika-fiber/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type ChangePasswordRequest struct {
	UserID          int    `json:"user_id"`
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

func ChangePassword(c *fiber.Ctx) error {
	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Response{Success: 0, Message: "Invalid request body", Data: nil})
	}

	if req.NewPassword != req.ConfirmPassword {
		return c.Status(http.StatusBadRequest).JSON(Response{Success: 0, Message: "New password and confirmation password do not match", Data: nil})
	}

	user, err := getUserID(req.UserID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{Success: 0, Message: "User not found", Data: nil})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(Response{Success: 0, Message: "Old password is incorrect", Data: nil})
	}

	hashedNewPassword, err := hasingPassword(req.NewPassword)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{Success: 0, Message: "Failed to hash new password", Data: nil})
	}

	err = updatePasswordInDatabase(req.UserID, hashedNewPassword)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{Success: 0, Message: "Failed to update password", Data: nil})
	}

	return c.Status(http.StatusOK).JSON(Response{Success: 1, Message: "Password changed successfully", Data: nil})
}

func getUserID(userID int) (models.User, error) {
	db := database.GetGormDBInstance()

	var user models.User
	err := db.First(&user, userID).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func updatePasswordInDatabase(userID int, newPassword string) error {
	db := database.GetGormDBInstance()

	err := db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password":   newPassword,
		"updated_at": time.Now(),
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func hasingPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
