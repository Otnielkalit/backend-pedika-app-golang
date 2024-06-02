package handlers

import (
	"net/http"
	"time"

	"backend-pedika-fiber/auth"
	"backend-pedika-fiber/database"
	"backend-pedika-fiber/helper"
	"backend-pedika-fiber/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetUserProfile(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	userID, err := auth.ExtractUserIDFromToken(tokenString)
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusUnauthorized,
			Status:  "error",
			Message: "Unauthorized",
		}
		return c.Status(http.StatusUnauthorized).JSON(response)
	}
	var user models.User
	if err := database.GetGormDBInstance().First(&user, userID).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to retrieve user profile",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "User profile retrieved successfully",
		Data:    user,
	}
	return c.Status(http.StatusOK).JSON(response)
}

func checkUsernameExists(db *gorm.DB, username string) bool {
	var count int64
	db.Model(&models.User{}).Where("username = ?", username).Count(&count)
	return count > 0
}

func checkEmailExists(db *gorm.DB, email string) bool {
	var count int64
	db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func checkPhoneNumberExists(db *gorm.DB, phoneNumber string) bool {
	var count int64
	db.Model(&models.User{}).Where("phone_number = ?", phoneNumber).Count(&count)
	return count > 0
}

func UpdateUserProfile(c *fiber.Ctx) error {
	var updateUser models.User
	if err := c.BodyParser(&updateUser); err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	userID, err := auth.ExtractUserIDFromToken(c.Get("Authorization"))
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to get user ID",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	tx := database.GetGormDBInstance().Begin()

	var existingUser models.User
	if err := tx.First(&existingUser, userID).Error; err != nil {
		tx.Rollback()
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to find user",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	if updateUser.Username != "" && updateUser.Username != existingUser.Username {
		if checkUsernameExists(tx, updateUser.Username) {
			tx.Rollback()
			response := helper.ResponseWithOutData{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: "Username ini sudah ada, coba yang lain",
			}
			return c.Status(http.StatusBadRequest).JSON(response)
		}
		existingUser.Username = updateUser.Username
	}

	if updateUser.Email != "" && updateUser.Email != existingUser.Email {
		if checkEmailExists(tx, updateUser.Email) {
			tx.Rollback()
			response := helper.ResponseWithOutData{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: "Email yang anda masukkan sudah pernah terdaftar",
			}
			return c.Status(http.StatusBadRequest).JSON(response)
		}
		existingUser.Email = updateUser.Email
	}

	if updateUser.PhoneNumber != "" && updateUser.PhoneNumber != existingUser.PhoneNumber {
		if checkPhoneNumberExists(tx, updateUser.PhoneNumber) {
			tx.Rollback()
			response := helper.ResponseWithOutData{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: "Nomor telepon yang anda masukkan sudah pernah terdaftar",
			}
			return c.Status(http.StatusBadRequest).JSON(response)
		}
		existingUser.PhoneNumber = updateUser.PhoneNumber
	}

	if updateUser.NIK != 0 && updateUser.NIK != existingUser.NIK {
		existingUser.NIK = updateUser.NIK
	}

	file, err := c.FormFile("photo_profile")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			tx.Rollback()
			response := helper.ResponseWithOutData{
				Code:    http.StatusInternalServerError,
				Status:  "error",
				Message: "Failed to open photo profile",
			}
			return c.Status(http.StatusInternalServerError).JSON(response)
		}
		defer src.Close()

		imageURL, err := helper.UploadFileToCloudinary(src, file.Filename)
		if err != nil {
			tx.Rollback()
			response := helper.ResponseWithOutData{
				Code:    http.StatusInternalServerError,
				Status:  "error",
				Message: "Failed to upload photo profile",
			}
			return c.Status(http.StatusInternalServerError).JSON(response)
		}
		existingUser.PhotoProfile = imageURL
	}

	tanggalLahirStr := c.FormValue("tanggal_lahir")
	if tanggalLahirStr != "" {
		tanggalLahir, err := time.Parse("02-01-2006", tanggalLahirStr)
		if err != nil {
			tx.Rollback()
			response := helper.ResponseWithOutData{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: "Invalid date format, use dd-MM-yyyy",
			}
			return c.Status(http.StatusBadRequest).JSON(response)
		}
		existingUser.TanggalLahir = tanggalLahir
	}

	if updateUser.Alamat != "" {
		existingUser.Alamat = updateUser.Alamat
	}

	if updateUser.FullName != "" {
		existingUser.FullName = updateUser.FullName
	}

	if updateUser.TempatLahir != "" {
		existingUser.TempatLahir = updateUser.TempatLahir
	}

	if updateUser.JenisKelamin != "" {
		existingUser.JenisKelamin = updateUser.JenisKelamin
	}

	existingUser.UpdatedAt = time.Now()

	if err := tx.Save(&existingUser).Error; err != nil {
		tx.Rollback()
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to update user profile",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	tx.Commit()

	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Profil Anda berhasil diupdate",
		Data:    existingUser,
	}
	return c.Status(http.StatusOK).JSON(response)
}
