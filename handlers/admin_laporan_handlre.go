package handlers

import (
	"backend-pedika-fiber/database"
	"backend-pedika-fiber/helper"
	"backend-pedika-fiber/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

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

func GetReportDetailsByID(c *fiber.Ctx) error {
	reportNoRegistrasi := c.Params("no_registrasi")

	type ReportDetails struct {
		NoRegistrasi        string           `json:"no_registrasi"`
		User                models.User      `json:"user_id" gorm:"foreignKey:UserID"`
		ViolenceCategoryID  int64            `json:"id_violence_category"`
		TanggalPelaporan    time.Time        `json:"tanggal_pelaporan"`
		TanggalKejadian     time.Time        `json:"tanggal_kejadian"`
		KategoriLokasiKasus string           `json:"kategori_lokasi_kasus"`
		Alamat              models.AlamatTKP `json:"id_alamat_tkp" gorm:"foreignKey:AlamatID"`
		KronologisKasus     string           `json:"kronologis_kasus"`
		CreatedAt           time.Time        `json:"created_at"`
		UpdatedAt           time.Time        `json:"updated_at"`
	}

	var reportDetails ReportDetails
	if err := database.GetGormDBInstance().
		Table("laporans").
		Select("laporans.*, users.*, alamat_tkps.*").
		Joins("JOIN users ON laporans.user_id = users.id").
		Joins("JOIN alamat_tkps ON laporans.id_alamat_tkp = alamat_tkps.id").
		Where("laporans.no_registrasi = ?", reportNoRegistrasi).
		First(&reportDetails).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusNotFound,
			Status:  "error",
			Message: "Report not found",
		}
		return c.Status(http.StatusNotFound).JSON(response)
	}
	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Report details retrieved successfully",
		Data:    reportDetails,
	}
	return c.Status(http.StatusOK).JSON(response)
}
