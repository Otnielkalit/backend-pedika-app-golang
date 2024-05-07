package models

import (
	"time"
)

type Korban struct {
	ID                   uint         `gorm:"primaryKey" json:"id"`
	NoRegistrasi         string       `json:"no_registrasi"`
	NIKPelaku            string       `json:"nik_pelaku"`
	Nama                 string       `json:"nama_pelaku"`
	Usia                 int          `json:"usia_pelaku"`
	AlamatPelaku         AlamatPelaku `gorm:"foreignKey:AlamatPelakuID"`
	AlamatPelakuID       uint         `json:"alamat_pelaku_id"`
	JenisKelamin         string       `json:"jenis_kelamin"`
	Agama                string       `json:"agama"`
	NoTelepon            string       `json:"no_telepon"`
	Pendidikan           string       `json:"pendidikan"`
	Pekerjaan            string       `json:"pekerjaan"`
	StatusPerkawinan     string       `json:"status_perkawinan"`
	Kebangsaan           string       `json:"kebangsaan"`
	HubunganDenganKorban string       `json:"hubungan_dengan_korban"`
	KeteranganLainnya    string       `json:"keterangan_lainnya"`
	DokumentasiPelaku    string       `json:"dokumentasi_korban"`
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
}
