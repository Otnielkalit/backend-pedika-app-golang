package migration

import (
	"backend-pedika-fiber/database"
	"backend-pedika-fiber/models"
	"log"
)

func RunMigration() {
	err := database.DB.AutoMigrate(
		&models.User{},
		&models.ViolenceCategory{},
		&models.EmergencyContact{},
		&models.Content{},
		&models.Laporan{},
		&models.AlamatTKP{},
		&models.Event{},
		&models.AlamatPelaku{},
		&models.Pelaku{})
	if err != nil {
		log.Println(err)
	}

}
