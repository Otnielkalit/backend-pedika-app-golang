package handlers

import (
	"net/http"
	"time"

	"backend-pedika-fiber/database"
	"backend-pedika-fiber/helper"
	"backend-pedika-fiber/models"

	"github.com/gofiber/fiber/v2"
)

func CreatePelaku(c *fiber.Ctx) error {
	var pelaku models.Pelaku
	if err := c.BodyParser(&pelaku); err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Failed to parse request body",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}
	file, err := c.FormFile("dokumentasi_pelaku")
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Failed to get dokumentasi pelaku",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	fileContent, err := file.Open()
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to open uploaded file",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	defer fileContent.Close()
	dokumentasiURL, err := helper.UploadFileToCloudinary(fileContent, file.Filename)
	if err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to upload dokumentasi pelaku",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	pelaku.DokumentasiPelaku = dokumentasiURL
	pelaku.NoRegistrasi = c.FormValue("no_registrasi")
	pelaku.NIKPelaku = c.FormValue("nik_pelaku")
	pelaku.JenisKelamin = c.FormValue("jenis_kelamin")
	pelaku.NoTelepon = c.FormValue("no_telepon")
	pelaku.StatusPerkawinan = c.FormValue("status_perkawinan")
	pelaku.HubunganDenganKorban = c.FormValue("hubungan_dengan_korban")
	pelaku.KeteranganLainnya = c.FormValue("keterangan_lainnya")
	now := time.Now()
	pelaku.CreatedAt = now
	pelaku.UpdatedAt = now

	alamatPelaku := models.AlamatPelaku{
		Provinsi:     c.FormValue("provinsi"),
		Kabupaten:    c.FormValue("kabupaten"),
		Kecamatan:    c.FormValue("kecamatan"),
		Desa:         c.FormValue("desa"),
		AlamatDetail: c.FormValue("alamat_detail"),
	}
	if err := database.GetGormDBInstance().Create(&alamatPelaku).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to create alamat pelaku",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	pelaku.AlamatPelakuID = alamatPelaku.ID
	if err := database.GetGormDBInstance().Create(&pelaku).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to create pelaku",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	responseData := map[string]interface{}{
		"id":                     pelaku.ID,
		"no_registrasi":          pelaku.NoRegistrasi,
		"nik_pelaku":             pelaku.NIKPelaku,
		"nama_pelaku":            pelaku.Nama,
		"usia_pelaku":            pelaku.Usia,
		"alamat_pelaku_id":       pelaku.AlamatPelakuID,
		"jenis_kelamin":          pelaku.JenisKelamin,
		"agama":                  pelaku.Agama,
		"no_telepon":             pelaku.NoTelepon,
		"pendidikan":             pelaku.Pendidikan,
		"pekerjaan":              pelaku.Pekerjaan,
		"status_perkawinan":      pelaku.StatusPerkawinan,
		"kebangsaan":             pelaku.Kebangsaan,
		"hubungan_dengan_korban": pelaku.HubunganDenganKorban,
		"keterangan_lainnya":     pelaku.KeteranganLainnya,
		"dokumentasi_korban":     pelaku.DokumentasiPelaku,
		"created_at":             pelaku.CreatedAt,
		"updated_at":             pelaku.UpdatedAt,
	}

	response := helper.ResponseWithData{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Pelaku created successfully",
		Data:    responseData,
	}
	return c.Status(http.StatusCreated).JSON(response)
}

func UpdatePelaku(c *fiber.Ctx) error {

	id := c.Params("id")
	var updatedPelaku models.Pelaku
	if err := database.GetGormDBInstance().First(&updatedPelaku, id).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusNotFound,
			Status:  "error",
			Message: "Pelaku not found",
		}
		return c.Status(http.StatusNotFound).JSON(response)
	}
	if err := c.BodyParser(&updatedPelaku); err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Failed to parse request body",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}
	updatedAlamatPelaku := models.AlamatPelaku{
		Provinsi:     c.FormValue("provinsi"),
		Kabupaten:    c.FormValue("kabupaten"),
		Kecamatan:    c.FormValue("kecamatan"),
		Desa:         c.FormValue("desa"),
		AlamatDetail: c.FormValue("alamat_detail"),
	}
	if err := database.GetGormDBInstance().Model(&updatedPelaku).Association("AlamatPelaku").Replace(&updatedAlamatPelaku); err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to update pelaku's address",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	if err := database.GetGormDBInstance().Save(&updatedPelaku).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to update pelaku",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	responseData := map[string]interface{}{
		"id":                     updatedPelaku.ID,
		"no_registrasi":          updatedPelaku.NoRegistrasi,
		"nik_pelaku":             updatedPelaku.NIKPelaku,
		"nama_pelaku":            updatedPelaku.Nama,
		"usia_pelaku":            updatedPelaku.Usia,
		"alamat_pelaku_id":       updatedPelaku.AlamatPelakuID,
		"jenis_kelamin":          updatedPelaku.JenisKelamin,
		"agama":                  updatedPelaku.Agama,
		"no_telepon":             updatedPelaku.NoTelepon,
		"pendidikan":             updatedPelaku.Pendidikan,
		"pekerjaan":              updatedPelaku.Pekerjaan,
		"status_perkawinan":      updatedPelaku.StatusPerkawinan,
		"kebangsaan":             updatedPelaku.Kebangsaan,
		"hubungan_dengan_korban": updatedPelaku.HubunganDenganKorban,
		"keterangan_lainnya":     updatedPelaku.KeteranganLainnya,
		"dokumentasi_korban":     updatedPelaku.DokumentasiPelaku,
		"created_at":             updatedPelaku.CreatedAt,
		"updated_at":             updatedPelaku.UpdatedAt,
	}

	response := helper.ResponseWithData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Pelaku updated successfully",
		Data:    responseData,
	}
	return c.Status(http.StatusOK).JSON(response)
}

func DeletePelaku(c *fiber.Ctx) error {
	id := c.Params("id")
	var pelaku models.Pelaku
	if err := database.GetGormDBInstance().First(&pelaku, id).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusNotFound,
			Status:  "error",
			Message: "Pelaku not found",
		}
		return c.Status(http.StatusNotFound).JSON(response)
	}

	if err := database.GetGormDBInstance().Delete(&pelaku).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to delete pelaku",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	if err := database.GetGormDBInstance().Where("id = ?", pelaku.AlamatPelakuID).Delete(&models.AlamatPelaku{}).Error; err != nil {
		response := helper.ResponseWithOutData{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to delete pelaku's address",
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := helper.ResponseWithOutData{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Pelaku deleted successfully",
	}
	return c.Status(http.StatusOK).JSON(response)
}
