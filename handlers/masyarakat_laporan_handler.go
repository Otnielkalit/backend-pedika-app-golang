package handlers

import (
	"backend-pedika-fiber/auth"
	"backend-pedika-fiber/database"
	"backend-pedika-fiber/helper"
	"backend-pedika-fiber/models"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

/*=========================== USER CREATE LAPORAN =======================*/
var mu sync.Mutex

func CreateLaporan(c *fiber.Ctx) error {
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
	var laporan models.Laporan
	if err := c.BodyParser(&laporan); err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}
	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("Error retrieving multipart form: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve multipart form",
		})
	}
	files := form.File["dokumentasi"]
	imageURLs, err := helper.UploadMultipleFileToCloudinary(files)
	if err != nil {
		log.Printf("Error uploading images to Cloudinary: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload images",
		})
	}
	laporan.Dokumentasi = datatypes.JSONMap{"urls": imageURLs}
	year := time.Now().Year()
	month := int(time.Now().Month())
	noRegistrasi, err := generateUniqueNoRegistrasi(month, year)
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to generate registration number",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	alamatTKP := models.AlamatTKP{
		Provinsi:     c.FormValue("provinsi"),
		Kabupaten:    c.FormValue("kabupaten"),
		Kecamatan:    c.FormValue("kecamatan"),
		Desa:         c.FormValue("desa"),
		AlamatDetail: c.FormValue("alamat_detail"),
	}
	if err := database.GetGormDBInstance().Create(&alamatTKP).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to create alamat TKP",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	laporan.NoRegistrasi = noRegistrasi
	laporan.TanggalPelaporan = time.Now()
	laporan.IDAlamatTKP = alamatTKP.ID
	laporan.UserID = uint(userID)
	laporan.CreatedAt = time.Now()
	laporan.UpdatedAt = time.Now()
	laporan.KategoriKekerasan = c.FormValue("kategori_kekerasan")
	laporan.KategoriLokasiKasus = c.FormValue("kategori_lokasi_kasus")
	laporan.KronologisKasus = c.FormValue("kronologis_kasus")
	laporan.TanggalKejadian = time.Time{}

	if err := database.GetGormDBInstance().Create(&laporan).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to create laporan",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	response := helper.ResponseWithData{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Laporan created successfully",
		Data: fiber.Map{
			"no_registrasi":         laporan.NoRegistrasi,
			"user_id":               laporan.UserID,
			"kategori_kekerasan":    laporan.KategoriKekerasan,
			"tanggal_pelaporan":     laporan.TanggalPelaporan,
			"tanggal_kejadian":      laporan.TanggalKejadian,
			"kategori_lokasi_kasus": laporan.KategoriLokasiKasus,
			"id_alamat_tkp":         laporan.IDAlamatTKP,
			"kronologis_kasus":      laporan.KronologisKasus,
			"dokumentasi": fiber.Map{
				"urls": imageURLs,
			},
			"created_at": laporan.CreatedAt,
			"updated_at": laporan.UpdatedAt,
		},
	}
	return c.Status(http.StatusCreated).JSON(response)
}

func generateUniqueNoRegistrasi(month, year int) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	romanMonth := convertToRoman(month)
	regNo := "001-DPMDPPA-" + romanMonth + "-" + strconv.Itoa(year)
	var existingCount int64
	if err := database.GetGormDBInstance().Model(&models.Laporan{}).Where("no_registrasi = ?", regNo).Count(&existingCount).Error; err != nil {
		return "", err
	}
	if existingCount > 0 {
		for i := 1; i < 1000; i++ {
			modifiedRegNo := fmt.Sprintf("%03d", i) + "-DPMDPPA-" + romanMonth + "-" + strconv.Itoa(year)
			var existingCount int64
			if err := database.GetGormDBInstance().Model(&models.Laporan{}).Where("no_registrasi = ?", modifiedRegNo).Count(&existingCount).Error; err != nil {
				return "", err
			}
			if existingCount == 0 {
				return modifiedRegNo, nil
			}
		}
		return "", errors.New("failed to generate unique registration number")
	}
	return regNo, nil
}
func convertToRoman(month int) string {
	months := [...]string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII"}
	if month >= 1 && month <= 12 {
		return months[month-1]
	}
	return ""
}

/*=========================== AMBIL SEMUA  LAPORAN SETIAP BERDASARKAN USER YANG LOGIN=======================*/
func GetUserReports(c *fiber.Ctx) error {
	userID, err := auth.ExtractUserIDFromToken(c.Get("Authorization"))
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusUnauthorized,
			Status:  "error",
			Message: "Unauthorized",
		}
		return c.Status(http.StatusUnauthorized).JSON(response)
	}
	reports, err := GetReportsByUserID(userID)
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to get user reports",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "List of laporan by user",
		Data:    reports,
	}
	return c.Status(http.StatusOK).JSON(response)
}

func GetReportsByUserID(userID uint) ([]map[string]interface{}, error) {
	var reports []models.Laporan
	if err := database.GetGormDBInstance().
		Where("user_id = ?", userID).
		Find(&reports).Error; err != nil {
		return nil, err
	}

	var formattedReports []map[string]interface{}
	for _, report := range reports {
		formattedReport := map[string]interface{}{
			"no_registrasi":         report.NoRegistrasi,
			"user_id":               report.UserID,
			"kategori_kekerasan":    report.KategoriKekerasan,
			"tanggal_pelaporan":     report.TanggalPelaporan,
			"tanggal_kejadian":      report.TanggalKejadian,
			"kategori_lokasi_kasus": report.KategoriLokasiKasus,
			"id_alamat_tkp":         report.IDAlamatTKP,
			"kronologis_kasus":      report.KronologisKasus,
			"dokumentasi":           report.Dokumentasi,
			"created_at":            report.CreatedAt,
			"updated_at":            report.UpdatedAt,
		}
		formattedReports = append(formattedReports, formattedReport)
	}

	return formattedReports, nil
}

/*=========================== TAMPILKAN DETAIL LAPORAN USER BERDASARKAN NO_REGISTRASI =======================*/
func GetReportByNoRegistrasi(c *fiber.Ctx) error {
	noRegistrasi := c.Params("no_registrasi")
	var laporan models.Laporan
	if err := database.GetGormDBInstance().
		Where("laporans.no_registrasi = ?", noRegistrasi).
		Preload("AlamatTKP").
		Preload("User").
		First(&laporan).Error; err != nil {
		status := http.StatusInternalServerError
		message := "Failed to fetch report detail"
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
			message = "Report not found"
		}
		response := helper.ResponseWithOutData{
			Code:    status,
			Status:  "error",
			Message: message,
		}
		return c.Status(status).JSON(response)
	}
	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Report detail retrieved successfully",
		Data:    laporan,
	}
	return c.Status(http.StatusOK).JSON(response)
}
