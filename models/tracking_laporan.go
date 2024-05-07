package models

import "gorm.io/datatypes"

type TrackingLaporan struct {
	ID           uint              `gorm:"primaryKey" json:"id"`
	Laporan      Laporan           `json:"NoRegistrasi"`
	NoRegistrasi string            `json:"no_registrasi"`
	Status       string            `json:"status"`
	Keterangan   string            `json:"keterangan"`
	Dokumentasi  datatypes.JSONMap `json:"dokumentasi" form:"image" gorm:"type:json"`
}
