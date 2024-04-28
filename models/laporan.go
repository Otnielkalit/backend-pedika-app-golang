package models

import (
	"time"

	"gorm.io/datatypes"
)

type Laporan struct {
    NoRegistrasi        string            `gorm:"primaryKey" json:"no_registrasi"`
    User                User              `gorm:"foreignKey:UserID"`
    UserID              uint              `json:"user_id"`
    KategoriKekerasan   string            `json:"kategori_kekerasan"`
    TanggalPelaporan    time.Time         `json:"tanggal_pelaporan"`
    TanggalKejadian     time.Time         `json:"tanggal_kejadian"`
    KategoriLokasiKasus string            `json:"kategori_lokasi_kasus"`
    AlamatTKP           AlamatTKP         `gorm:"foreignKey:IDAlamatTKP"`
    IDAlamatTKP         uint              `json:"id_alamat_tkp"`
    KronologisKasus     string            `json:"kronologis_kasus"`
    Dokumentasi         datatypes.JSONMap `json:"dokumentasi" form:"image" gorm:"type:json"`
    CreatedAt           time.Time         `json:"created_at"`
    UpdatedAt           time.Time         `json:"updated_at"`
}
 