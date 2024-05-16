package handlers

import (
	"backend-pedika-fiber/database"
	"backend-pedika-fiber/helper"
	"backend-pedika-fiber/models"
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func CreateTrackingLaporan(c *fiber.Ctx) error {
	var trackingLaporan models.TrackingLaporan
	if err := c.BodyParser(&trackingLaporan); err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	// Explicitly fetch no_registrasi from form data
	noRegistrasi := c.FormValue("no_registrasi")
	if noRegistrasi == "" {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "No Registrasi is required",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}
	var existingLaporan models.Laporan
	if err := database.GetGormDBInstance().Where("no_registrasi = ?", noRegistrasi).First(&existingLaporan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := helper.ResponseWithOutData{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: "No Registrasi not found in Laporan table",
			}
			return c.Status(http.StatusBadRequest).JSON(response)
		}
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Database error",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve multipart form",
		})
	}
	files := form.File["document"]
	imageURLs, err := helper.UploadMultipleFileToCloudinary(files)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload documents",
		})
	}

	trackingLaporan.Document = datatypes.JSONMap{"urls": imageURLs}
	trackingLaporan.NoRegistrasi = noRegistrasi
	trackingLaporan.Keterangan = c.FormValue("keterangan")
	trackingLaporan.CreatedAt = time.Now()
	trackingLaporan.UpdatedAt = time.Now()

	if err := database.GetGormDBInstance().Create(&trackingLaporan).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to create tracking laporan",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := helper.ResponseWithData{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Tracking laporan created successfully",
		Data: fiber.Map{
			"id":            trackingLaporan.ID,
			"no_registrasi": trackingLaporan.NoRegistrasi,
			"keterangan":    trackingLaporan.Keterangan,
			"document":      trackingLaporan.Document,
			"created_at":    trackingLaporan.CreatedAt,
			"updated_at":    trackingLaporan.UpdatedAt,
		},
	}

	return c.Status(http.StatusCreated).JSON(response)
}

func UpdateTrackingLaporan(c *fiber.Ctx) error {
	trackingLaporanID := c.Params("id")
	if trackingLaporanID == "" {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "ID is required",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	var trackingLaporan models.TrackingLaporan
	if err := database.GetGormDBInstance().First(&trackingLaporan, trackingLaporanID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := helper.ResponseWithOutData{
				Code:    http.StatusNotFound,
				Status:  "error",
				Message: "Tracking Laporan not found",
			}
			return c.Status(http.StatusNotFound).JSON(response)
		}
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Database error",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	var updatedData models.TrackingLaporan
	if err := c.BodyParser(&updatedData); err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	// Update only the fields that are not empty in the request body
	if updatedData.NoRegistrasi != "" {
		trackingLaporan.NoRegistrasi = updatedData.NoRegistrasi
	}
	if updatedData.Keterangan != "" {
		trackingLaporan.Keterangan = updatedData.Keterangan
	}

	// Handle file uploads
	form, err := c.MultipartForm()
	if err != nil && err != http.ErrNotMultipart {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve multipart form",
		})
	}

	if form != nil {
		files := form.File["document"]
		if len(files) > 0 {
			imageURLs, err := helper.UploadMultipleFileToCloudinary(files)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to upload images",
				})
			}
			trackingLaporan.Document = datatypes.JSONMap{"urls": imageURLs}
		}
	} else if updatedData.Document != nil {
		trackingLaporan.Document = updatedData.Document
	}

	trackingLaporan.UpdatedAt = time.Now()

	if err := database.GetGormDBInstance().Save(&trackingLaporan).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to update tracking laporan",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Tracking laporan updated successfully",
		Data: fiber.Map{
			"id":            trackingLaporan.ID,
			"no_registrasi": trackingLaporan.NoRegistrasi,
			"keterangan":    trackingLaporan.Keterangan,
			"document":      trackingLaporan.Document,
			"created_at":    trackingLaporan.CreatedAt,
			"updated_at":    trackingLaporan.UpdatedAt,
		},
	}

	return c.Status(http.StatusOK).JSON(response)
}

func DeleteTrackingLaporan(c *fiber.Ctx) error {
	trackingLaporanID := c.Params("id")
	if trackingLaporanID == "" {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "ID is required",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	var trackingLaporan models.TrackingLaporan
	if err := database.GetGormDBInstance().First(&trackingLaporan, trackingLaporanID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := helper.ResponseWithOutData{
				Code:    http.StatusNotFound,
				Status:  "error",
				Message: "Tracking Laporan not found",
			}
			return c.Status(http.StatusNotFound).JSON(response)
		}
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Database error",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	if err := database.GetGormDBInstance().Delete(&trackingLaporan).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to delete tracking laporan",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := helper.ResponseWithOutData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Tracking laporan deleted successfully",
	}

	return c.Status(http.StatusOK).JSON(response)
}
