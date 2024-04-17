package handlers

import (
	"backend-pedika-fiber/auth"
	"backend-pedika-fiber/database"
	"backend-pedika-fiber/helper"
	"backend-pedika-fiber/models"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

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
		NoRegistrasi: noRegistrasi,
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

	idViolenceCategoryStr := c.FormValue("id_violence_category")

	idViolenceCategory, err := strconv.Atoi(idViolenceCategoryStr)
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid ID Violence Category",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}
	laporan.IDViolenceCategory = idViolenceCategory
	laporan.KategoriLokasiKasus = c.FormValue("kategori_lokasi_kasus")
	laporan.KronologisKasus = c.FormValue("kronologis_kasus")
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
		Data:    laporan,
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
		Message: "User reports retrieved successfully",
		Data:    reports,
	}
	return c.Status(http.StatusOK).JSON(response)
}


func GetReportsByUserID(userID uint) ([]models.Laporan, error) {
	var reports []models.Laporan
	if err := database.GetGormDBInstance().Where("user_id = ?", userID).Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}
