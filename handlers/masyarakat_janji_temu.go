package handlers

import (
	"backend-pedika-fiber/auth"
	"backend-pedika-fiber/database"
	"backend-pedika-fiber/helper"
	"backend-pedika-fiber/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func MasyarakatCreateJanjiTemu(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	userID, err := auth.ExtractUserIDFromToken(token)
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusUnauthorized,
			Status:  "error",
			Message: "Unauthorized",
		}
		return c.Status(http.StatusUnauthorized).JSON(response)
	}

	var janjitemu models.JanjiTemu
	if err := c.BodyParser(&janjitemu); err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}
	waktuDimulai, err := time.Parse("2006-01-02T15:04:05", c.FormValue("waktu_dimulai"))
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid format for start time",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}
	waktuSelesai, err := time.Parse("2006-01-02T15:04:05", c.FormValue("waktu_selesai"))
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid format for end time",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}
	janjitemu.WaktuDimulai = waktuDimulai
	janjitemu.WaktuSelesai = waktuSelesai
	janjitemu.Status = "Belum disetujui"
	janjitemu.KeperluanKonsultasi = c.FormValue("keperluan_konsultasi")
	janjitemu.UserID = uint(userID)
	janjitemu.UserIDTolakSetujui = nil

	if err := database.DB.Create(&janjitemu).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to create janjitemu",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	responseData := struct {
		ID                  uint      `json:"id"`
		UserID              uint      `json:"user_id"`
		WaktuDimulai        time.Time `json:"waktu_dimulai"`
		WaktuSelesai        time.Time `json:"waktu_selesai"`
		KeperluanKonsultasi string    `json:"keperluan_konsultasi"`
		Status              string    `json:"status"`
		UserTolakSetujui    uint      `json:"user_tolak_setujui"`
		AlasanDitolak       string    `json:"alasan_ditolak"`
		AlasanDibatalkan    string    `json:"alasan_dibatalkan"`
	}{
		ID:                  janjitemu.ID,
		UserID:              janjitemu.UserID,
		WaktuDimulai:        janjitemu.WaktuDimulai,
		WaktuSelesai:        janjitemu.WaktuSelesai,
		KeperluanKonsultasi: janjitemu.KeperluanKonsultasi,
		Status:              janjitemu.Status,
		UserTolakSetujui:    0,
		AlasanDitolak:       janjitemu.AlasanDitolak,
		AlasanDibatalkan:    janjitemu.AlasanDibatalkan,
	}

	response := helper.ResponseWithData{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Janjitemu created successfully",
		Data:    responseData,
	}
	return c.Status(http.StatusCreated).JSON(response)
}

func MasyarakatEditJanjiTemu(c *fiber.Ctx) error {
	janjiTemuID := c.Params("id")

	var updateRequest struct {
		WaktuDimulai        time.Time `json:"waktu_dimulai"`
		WaktuSelesai        time.Time `json:"waktu_selesai"`
		KeperluanKonsultasi string    `json:"keperluan_konsultasi"`
	}
	if err := c.BodyParser(&updateRequest); err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	waktuDimulai := updateRequest.WaktuDimulai
	waktuSelesai := updateRequest.WaktuSelesai

	var janjiTemu models.JanjiTemu
	if err := database.DB.First(&janjiTemu, janjiTemuID).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusNotFound,
			Status:  "error",
			Message: "Janji temu not found",
		}
		return c.Status(http.StatusNotFound).JSON(response)
	}
	if janjiTemu.Status != "Belum disetujui" {
		response := helper.ResponseWithOutData{
			Code:    http.StatusForbidden,
			Status:  "error",
			Message: "Forbidden: You can only edit appointments with status 'Belum disetujui'",
		}
		return c.Status(http.StatusForbidden).JSON(response)
	}
	janjiTemu.WaktuDimulai = waktuDimulai
	janjiTemu.WaktuSelesai = waktuSelesai
	janjiTemu.KeperluanKonsultasi = c.FormValue("keperluan_konsultasi")

	if err := database.DB.Save(&janjiTemu).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to update janji temu",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := helper.ResponseWithOutData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Janji temu updated successfully",
	}
	return c.Status(http.StatusOK).JSON(response)
}

func GetUserJanjiTemus(c *fiber.Ctx) error {
	userID, err := auth.ExtractUserIDFromToken(c.Get("Authorization"))
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusUnauthorized,
			Status:  "error",
			Message: "Unauthorized",
		}
		return c.Status(http.StatusUnauthorized).JSON(response)
	}
	var janjiTemus []models.JanjiTemu
	if err := database.DB.Preload("User").Preload("UserTolakSetujui").Where("user_id = ?", userID).Find(&janjiTemus).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to get user JanjiTemu records",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	if len(janjiTemus) == 0 {
		response := helper.ResponseWithOutData{
			Code:    http.StatusOK,
			Status:  "success",
			Message: "No JanjiTemu records found for the user",
		}
		return c.Status(http.StatusOK).JSON(response)
	}
	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "List of JanjiTemu by user",
		Data:    janjiTemus,
	}
	return c.Status(http.StatusOK).JSON(response)
}

func GetJanjiTemuByID(c *fiber.Ctx) error {
	janjiTemuID := c.Params("id")
	var janjiTemu models.JanjiTemu
	if err := database.DB.Preload("UserTolakSetujui").First(&janjiTemu, janjiTemuID).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusNotFound,
			Status:  "error",
			Message: "JanjiTemu not found",
		}
		return c.Status(http.StatusNotFound).JSON(response)
	}
	if janjiTemu.Status == "Ditolak" && janjiTemu.UserTolakSetujui.ID != 0 {
		var user models.User
		if err := database.DB.First(&user, janjiTemu.UserIDTolakSetujui).Error; err != nil {
			response := helper.ResponseWithOutData{
				Code:    http.StatusInternalServerError,
				Status:  "error",
				Message: "Failed to fetch user detail who rejected the appointment",
			}
			return c.Status(http.StatusInternalServerError).JSON(response)
		}
		janjiTemu.UserTolakSetujui = user
	}
	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "JanjiTemu detail",
		Data:    janjiTemu,
	}
	return c.Status(http.StatusOK).JSON(response)
}

func MasyarakatCancelJanjiTemu(c *fiber.Ctx) error {
	janjiTemuID := c.Params("id")

	var janjiTemu models.JanjiTemu
	if err := database.DB.First(&janjiTemu, janjiTemuID).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusNotFound,
			Status:  "error",
			Message: "Janji temu not found",
		}
		return c.Status(http.StatusNotFound).JSON(response)
	}
	if janjiTemu.Status != "Belum disetujui" {
		response := helper.ResponseWithOutData{
			Code:    http.StatusForbidden,
			Status:  "error",
			Message: "Forbidden: You can only cancel appointments with status 'Belum disetujui'",
		}
		return c.Status(http.StatusForbidden).JSON(response)
	}
	janjiTemu.Status = "Dibatalkan"
	janjiTemu.AlasanDibatalkan = c.FormValue("alasan_dibatalkan")

	if err := database.DB.Save(&janjiTemu).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to cancel janji temu",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := helper.ResponseWithOutData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Janji temu canceled successfully",
	}
	return c.Status(http.StatusOK).JSON(response)
}

func AdminGetAllJanjiTemu(c *fiber.Ctx) error {
	var janjiTemus []models.JanjiTemu
	if err := database.DB.Preload("User").Preload("UserTolakSetujui").Find(&janjiTemus).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to retrieve Janji Temu data",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "List of Janji Temu",
		Data:    janjiTemus,
	}
	return c.Status(http.StatusOK).JSON(response)
}

func AdminJanjiTemuByID(c *fiber.Ctx) error {
	janjiTemuID := c.Params("id")
	var janjiTemu models.JanjiTemu
	if err := database.DB.Preload("UserTolakSetujui").Preload("User").First(&janjiTemu, janjiTemuID).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusNotFound,
			Status:  "error",
			Message: "JanjiTemu not found",
		}
		return c.Status(http.StatusNotFound).JSON(response)
	}
	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "JanjiTemu detail",
		Data:    janjiTemu,
	}
	return c.Status(http.StatusOK).JSON(response)
}

func AdminApproveJanjiTemu(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	userID, err := auth.ExtractUserIDFromToken(token)
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusUnauthorized,
			Status:  "error",
			Message: "Unauthorized",
		}
		return c.Status(http.StatusUnauthorized).JSON(response)
	}
	id := c.Params("id")
	var janjiTemu models.JanjiTemu
	if err := database.DB.First(&janjiTemu, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(helper.ResponseWithOutData{
			Code:    http.StatusNotFound,
			Status:  "error",
			Message: "Janji temu tidak ditemukan",
		})
	}
	janjiTemu.UserIDTolakSetujui = &userID
	janjiTemu.Status = "Disetujui"

	if err := database.DB.Save(&janjiTemu).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Gagal menyimpan perubahan status",
		})
	}
	response := helper.ResponseWithOutData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Status janji temu berhasil diubah menjadi Disetujui",
	}
	return c.Status(http.StatusOK).JSON(response)
}

func AdminCancelJanjiTemu(c *fiber.Ctx) error {
	janjiTemuID := c.Params("id")
	token := c.Get("Authorization")
	userID, err := auth.ExtractUserIDFromToken(token)
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusUnauthorized,
			Status:  "error",
			Message: "Unauthorized",
		}
		return c.Status(http.StatusUnauthorized).JSON(response)
	}
	var cancelRequest struct {
		AlasanDitolak string `json:"alasan_ditolak"`
	}
	if err := c.BodyParser(&cancelRequest); err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}
	var janjiTemu models.JanjiTemu
	if err := database.DB.First(&janjiTemu, janjiTemuID).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusNotFound,
			Status:  "error",
			Message: "Janji temu not found",
		}
		return c.Status(http.StatusNotFound).JSON(response)
	}
	janjiTemu.Status = "Ditolak"
	janjiTemu.UserIDTolakSetujui = &userID
	janjiTemu.AlasanDitolak = c.FormValue("alasan_ditolak")
	if err := database.DB.Save(&janjiTemu).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to cancel janji temu",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	response := helper.ResponseWithOutData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Janji temu canceled successfully",
	}
	return c.Status(http.StatusOK).JSON(response)
}
