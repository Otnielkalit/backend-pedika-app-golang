package handlers

import (
	"net/http"
	"strings"

	"backend-pedika-fiber/auth"
	"backend-pedika-fiber/database"
	"backend-pedika-fiber/helper"
	"backend-pedika-fiber/models"

	"github.com/gofiber/fiber/v2"
)

func GetUserProfile(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	userID, err := auth.ExtractUserIDFromToken(tokenString)
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusUnauthorized,
			Status:  "error",
			Message: "Unauthorized",
		}
		return c.Status(http.StatusUnauthorized).JSON(response)
	}
	var user models.User
	if err := database.GetGormDBInstance().First(&user, userID).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to retrieve user profile",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "User profile retrieved successfully",
		Data:    user,
	}
	return c.Status(http.StatusOK).JSON(response)
}

func EditProfile(c *fiber.Ctx) error {
	var updateUser models.User
	if err := c.BodyParser(&updateUser); err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}
	userID, err := auth.ExtractUserIDFromToken(c.Get("Authorization"))
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to get user ID",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	var existingUser models.User
	if err := database.GetGormDBInstance().First(&existingUser, userID).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to find user",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	if updateUser.PhotoProfile != "" {
		imageUrl, err := helper.UploadFileToCloudinary(strings.NewReader(updateUser.PhotoProfile), "profile_picture")
		if err != nil {
			response := helper.ResponseWithOutData{
				Code:    http.StatusInternalServerError,
				Status:  "error",
				Message: "Failed to upload profile picture",
			}
			return c.Status(http.StatusInternalServerError).JSON(response)
		}
		existingUser.PhotoProfile = imageUrl
	}
	existingUser.Username = updateUser.FullName
	existingUser.Username = updateUser.Username
	existingUser.Username = updateUser.PhoneNumber
	existingUser.Email = updateUser.Email
	if err := database.GetGormDBInstance().Save(&existingUser).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to update user profile",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "User profile updated successfully",
		Data:    existingUser,
	}
	return c.Status(http.StatusOK).JSON(response)
}
