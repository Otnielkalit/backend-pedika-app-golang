package handlers

import (
	"backend-pedika-fiber/database"
	"backend-pedika-fiber/helper"
	"backend-pedika-fiber/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetAllViolenceCategories(c *fiber.Ctx) error {
	var categories []models.ViolenceCategory
	if err := database.DB.Find(&categories).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Internal Server Error",
		})
	}
	return c.Status(http.StatusOK).JSON(helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "List of violence categories",
		Data:    categories,
	})
}

func GetViolenceCategoryByID(c *fiber.Ctx) error {
	categoryID := c.Params("id")

	var category models.ViolenceCategory
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(helper.ResponseWithOutData{
			Code:    http.StatusNotFound,
			Status:  "error",
			Message: "Violence category not found",
		})
	}

	return c.Status(http.StatusOK).JSON(helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Violence category details",
		Data:    category,
	})
}

func CreateViolenceCategory(c *fiber.Ctx) error {
	var category models.ViolenceCategory
	if err := c.BodyParser(&category); err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	file, err := c.FormFile("image")
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

	category.Image = imageURL
	category.CategoryName = c.FormValue("category_name")

	if err := database.DB.Create(&category).Error; err != nil {
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
		Data:    category,
	}
	return c.Status(http.StatusCreated).JSON(response)
}

func UpdateViolenceCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var category models.ViolenceCategory
	if err := database.DB.First(&category, id).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusNotFound,
			Status:  "error",
			Message: "Category not found",
		}
		return c.Status(http.StatusNotFound).JSON(response)
	}
	categoryName := c.FormValue("category_name")
	if categoryName != "" {
		category.CategoryName = categoryName
	}
	file, err := c.FormFile("image")
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

		category.Image = imageURL
	}
	category.UpdatedAt = time.Now()

	if err := database.DB.Save(&category).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to update violence category",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Violence category updated successfully",
		Data:    category,
	}
	return c.Status(http.StatusOK).JSON(response)
}

func DeleteViolenceCategory(c *fiber.Ctx) error {
	categoryID := c.Params("id")

	var deletedCategory models.ViolenceCategory
	if err := database.DB.First(&deletedCategory, categoryID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(helper.ResponseWithOutData{
			Code:    http.StatusNotFound,
			Status:  "error",
			Message: "Violence category not found",
		})
	}

	if err := database.DB.Where("id = ?", categoryID).Delete(&models.ViolenceCategory{}).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to delete violence category",
		})
	}

	return c.Status(http.StatusOK).JSON(helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Violence category deleted successfully",
	})
}
