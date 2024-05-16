package handlers

import (
	"backend-pedika-fiber/database"
	"backend-pedika-fiber/helper"
	"backend-pedika-fiber/models"
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/*=========================== AMBIL SEMUA LAPORAN =======================*/
func GetLatestReports(c *fiber.Ctx) error {
	var reports []models.Laporan
	db := database.GetGormDBInstance()

	if err := db.
		Preload("ViolenceCategory").
		Order("created_at desc").
		Limit(10).
		Find(&reports).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to fetch latest reports",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	var result []map[string]interface{}
	for _, report := range reports {
		result = append(result, map[string]interface{}{
			"no_registrasi":         report.NoRegistrasi,
			"user_id":               report.UserID,
			"violence_category":     report.ViolenceCategory,
			"kategori_kekerasan_id": report.KategoriKekerasanID,
			"tanggal_pelaporan":     report.TanggalPelaporan,
			"tanggal_kejadian":      report.TanggalKejadian,
			"kategori_lokasi_kasus": report.KategoriLokasiKasus,
			"alamat_tkp":            report.AlamatTKP,
			"alamat_detail_tkp":     report.AlamatDetailTKP,
			"kronologis_kasus":      report.KronologisKasus,
			"status":                report.Status,
			"alasan_dibatalkan":     report.AlasanDibatalkan,
			"waktu_dibatalkan":      report.WaktuDibatalkan,
			"waktu_dilihat":         report.WaktuDilihat,
			"waktu_diproses":        report.WaktuDiproses,
			"dokumentasi":           report.Dokumentasi,
			"created_at":            report.CreatedAt,
			"updated_at":            report.UpdatedAt,
		})
	}

	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Latest reports retrieved successfully",
		Data:    result,
	}
	return c.Status(http.StatusOK).JSON(response)
}

/*=========================== TAMPILKAN DETAIL LAPORAN USER BERDASARKAN NO_REGISTRASI =======================*/

func GetLaporanByNoRegistrasi(c *fiber.Ctx) error {
	noRegistrasi := c.Params("no_registrasi")
	var laporan models.Laporan

	// Get the laporan details with Preload
	if err := database.GetGormDBInstance().
		Preload("User").
		Preload("ViolenceCategory").
		Where("no_registrasi = ?", noRegistrasi).
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
	var trackingLaporan []models.TrackingLaporan
	if err := database.GetGormDBInstance().Where("no_registrasi = ?", noRegistrasi).Find(&trackingLaporan).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to fetch tracking laporan details",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	responseData := struct {
		models.Laporan
		TrackingLaporan []models.TrackingLaporan `json:"tracking_laporan"`
	}{
		Laporan:         laporan,
		TrackingLaporan: trackingLaporan,
	}

	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Report detail retrieved successfully",
		Data:    responseData,
	}
	return c.Status(http.StatusOK).JSON(response)
}

func AdminLihatLaporan(c *fiber.Ctx) error {
	noRegistrasi := c.Params("no_registrasi")

	var laporan models.Laporan
	db := database.GetGormDBInstance()
	if err := db.Where("no_registrasi = ?", noRegistrasi).First(&laporan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := helper.ResponseWithOutData{
				Code:    http.StatusNotFound,
				Status:  "error",
				Message: "Laporan not found",
			}
			return c.Status(http.StatusNotFound).JSON(response)
		}
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to retrieve laporan",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	laporan.Status = "Dilihat"
	now := time.Now()
	laporan.WaktuDilihat = &now

	if err := db.Save(&laporan).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to update laporan",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Laporan status updated successfully",
		Data: fiber.Map{
			"no_registrasi": laporan.NoRegistrasi,
			"status":        laporan.Status,
			"waktu_dilihat": laporan.WaktuDilihat,
			"updated_at":    laporan.UpdatedAt,
		},
	}

	return c.Status(http.StatusOK).JSON(response)
}

func AdminProsesLaporan(c *fiber.Ctx) error {
	noRegistrasi := c.Params("no_registrasi")

	var laporan models.Laporan
	db := database.GetGormDBInstance()
	if err := db.Where("no_registrasi = ?", noRegistrasi).First(&laporan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := helper.ResponseWithOutData{
				Code:    http.StatusNotFound,
				Status:  "error",
				Message: "Laporan not found",
			}
			return c.Status(http.StatusNotFound).JSON(response)
		}
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to retrieve laporan",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	laporan.Status = "Diproses"
	now := time.Now()
	laporan.WaktuDiproses = &now

	if err := db.Save(&laporan).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to update laporan",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Laporan status updated to Diproses successfully",
		Data: fiber.Map{
			"no_registrasi":  laporan.NoRegistrasi,
			"status":         laporan.Status,
			"waktu_diproses": laporan.WaktuDiproses,
			"updated_at":     laporan.UpdatedAt,
		},
	}

	return c.Status(http.StatusOK).JSON(response)
}
