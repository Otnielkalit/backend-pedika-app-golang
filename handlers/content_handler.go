package handlers

import (
	"backend-pedika-fiber/database"
	"backend-pedika-fiber/helper"
	"backend-pedika-fiber/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func GetAllContents(c *fiber.Ctx) error {
	var contents []models.Content
	if err := database.DB.Find(&contents).Error; err != nil {
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
		Data:    contents,
	}
	return c.Status(http.StatusOK).JSON(response)
}

func CreateContent(c *fiber.Ctx) error {
	var content models.Content
	if err := c.BodyParser(&content); err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	file, err := c.FormFile("image_content")
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

	content.ImageContent = imageURL
	content.Judul = c.FormValue("judul")
	content.IsiContent = c.FormValue("isi_content")

	if err := database.DB.Create(&content).Error; err != nil {
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
		Data:    content,
	}
	return c.Status(http.StatusCreated).JSON(response)
}

func UpdateContent(c *fiber.Ctx) error {
	contentID := c.Params("id")
	var existingContent models.Content
	if err := database.DB.First(&existingContent, contentID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Content not found",
		})
	}
	var updatedContent models.Content
	if err := c.BodyParser(&updatedContent); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if updatedContent.ImageContent != "" {
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
		updatedContent.ImageContent = imageURL
	}
	if err := database.DB.Model(&models.Content{}).Where("id = ?", contentID).Updates(&updatedContent).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update content",
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Content updated successfully",
	})
}

func DeleteContent(c *fiber.Ctx) error {
	contentID := c.Params("id")
	var content models.Content
	if err := database.DB.First(&content, contentID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Content not found",
		})
	}
	if err := database.DB.Delete(&content, contentID).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete content",
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Content deleted successfully",
	})
}
