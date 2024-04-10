package models

type AlamatTKP struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	NoRegistrasi string `json:"no_registrasi"`
	Provinsi     string `json:"provinsi"`
	Kabupaten    string `json:"kabupaten"`
	Kecamatan    string `json:"kecamatan"`
	Desa         string `json:"desa"`
	AlamatDetail string `json:"alamat_detail"`
}
