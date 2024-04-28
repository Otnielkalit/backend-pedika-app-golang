package handlers

import (
	"backend-pedika-fiber/database"
	"backend-pedika-fiber/helper"
	"backend-pedika-fiber/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetAllEvent(c *fiber.Ctx) error {
	var event []models.Event
	if err := database.DB.Find(&event).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to retrieve contents",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "List of contents",
		Data:    event,
	}
	return c.Status(http.StatusOK).JSON(response)
}

func CreateEvent(c *fiber.Ctx) error {
	var event models.Event
	if err := c.BodyParser(&event); err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	file, err := c.FormFile("thumbnail_event")
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Image file not provided",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

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

	event.ThumbnailEvent = imageURL
	event.NamaEvent = c.FormValue("nama_event")
	event.DeskripsiEvent = c.FormValue("deskripsi_event")
	tanggalPelaksanaanStr := c.FormValue("tanggal_pelaksanaan")
	tanggalPelaksanaan, err := time.Parse("02-01-2006", tanggalPelaksanaanStr) // Format tanggal hari-bulan-tahun
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid date format",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}
	event.TanggalPelaksanaan = tanggalPelaksanaan

	if err := database.DB.Create(&event).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to create violence category",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := helper.ResponseWithData{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Violence category created successfully",
		Data:    event,
	}
	return c.Status(http.StatusCreated).JSON(response)
}

func UpdateEvent(c *fiber.Ctx) error {
	eventID := c.Params("id")
	var existingContent models.Event
	if err := database.DB.First(&existingContent, eventID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Content not found",
		})
	}
	var updatedEvent models.Event
	if err := c.BodyParser(&updatedEvent); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if updatedEvent.ThumbnailEvent != "" {
		file, err := c.FormFile("image_content")
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Image file not provided",
			})
		}

		src, err := file.Open()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to open image file",
			})
		}
		defer src.Close()

		imageURL, err := helper.UploadFileToCloudinary(src, file.Filename)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to upload image",
			})
		}
		updatedEvent.ThumbnailEvent = imageURL
	}
	if err := database.DB.Model(&models.Event{}).Where("id = ?", eventID).Updates(&updatedEvent).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update content",
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Content updated successfully",
	})
}

func DeleteEvent(c *fiber.Ctx) error {
	eventID := c.Params("id")
	var event models.Event
	if err := database.DB.First(&event, eventID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "event not found",
		})
	}
	if err := database.DB.Delete(&event, eventID).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete event",
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "event deleted successfully",
	})
}
