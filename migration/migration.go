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
		&models.Korban{},
		&models.Pelaku{},
		&models.TrackingLaporan{},
		&models.Event{},
		&models.JanjiTemu{})
	if err != nil {
		log.Println(err)
	}

}
