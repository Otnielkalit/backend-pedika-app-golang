package handlers

import (
	"backend-pedika-fiber/database"
	"backend-pedika-fiber/helper"
	"backend-pedika-fiber/models"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/*=========================== AMBIL SEMUA LAPORAN =======================*/
func GetLatestReports(c *fiber.Ctx) error {
	var reports []models.Laporan
	if err := database.GetGormDBInstance().
		Table("laporans").
		Select("laporans.*, users.*, alamat_tkps.*").
		Joins("JOIN users ON laporans.user_id = users.id").
		Joins("JOIN alamat_tkps ON laporans.id_alamat_tkp = alamat_tkps.id").
		Order("laporans.created_at desc").
		Limit(10).
		Find(&reports).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to fetch latest reports",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Latest reports retrieved successfully",
		Data:    reports,
	}
	return c.Status(http.StatusOK).JSON(response)
}

/*=========================== TAMPILKAN DETAIL LAPORAN USER BERDASARKAN NO_REGISTRASI =======================*/
func GetLaporanByNoRegistrasi(c *fiber.Ctx) error {
	noRegistrasi := c.Params("no_registrasi")
	var laporan models.Laporan
	if err := database.GetGormDBInstance().
		Where("laporans.no_registrasi = ?", noRegistrasi).
		Preload("AlamatTKP").
		Preload("User").
		First(&laporan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := helper.ResponseWithOutData{
				Code:    http.StatusNotFound,
				Status:  "error",
				Message: "Report not found",
			}
			return c.Status(http.StatusNotFound).JSON(response)
		}
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to fetch report detail",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Report detail retrieved successfully",
		Data:    laporan,
	}
	return c.Status(http.StatusOK).JSON(response)
}
