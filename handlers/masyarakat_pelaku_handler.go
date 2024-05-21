package handlers

import (
	"backend-pedika-fiber/database"
	"backend-pedika-fiber/helper"
	"backend-pedika-fiber/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreatePelaku(c *fiber.Ctx) error {
	var pelaku models.Pelaku

	// Parse the request body to pelaku
	if err := c.BodyParser(&pelaku); err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	// Check if a new image file is provided
	file, err := c.FormFile("dokumentasi_pelaku")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			response := helper.ResponseWithOutData{
				Code:    http.StatusInternalServerError,
				Status:  "error",
				Message: "Failed to open image file",
			}
			return c.Status(http.StatusInternalServerError).JSON(response)
		}
		defer src.Close()

		imageURL, err := helper.UploadFileToCloudinary(src, file.Filename)
		if err != nil {
			response := helper.ResponseWithOutData{
				Code:    http.StatusInternalServerError,
				Status:  "error",
				Message: "Failed to upload image",
			}
			return c.Status(http.StatusInternalServerError).JSON(response)
		}

		pelaku.DokumentasiPelaku = imageURL
	}

	// Set the timestamps
	pelaku.CreatedAt = time.Now()
	pelaku.UpdatedAt = time.Now()

	// Save the new pelaku to the database
	if err := database.DB.Create(&pelaku).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to create pelaku",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	// Create success response
	response := helper.ResponseWithData{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Pelaku created successfully",
		Data:    pelaku,
	}
	return c.Status(http.StatusCreated).JSON(response)
}
