package models

import (
	"time"

	"gorm.io/datatypes"
)

type Laporan struct {
	NoRegistrasi        string            `gorm:"primaryKey" json:"no_registrasi"`
	User                User              `gorm:"foreignKey:UserID"`
	UserID              uint              `json:"user_id"`
	ViolenceCategory    ViolenceCategory  `gorm:"foreignKey:KategoriKekerasanID"`
	KategoriKekerasanID uint              `json:"kategori_kekerasan_id"`
	TanggalPelaporan    time.Time         `json:"tanggal_pelaporan"`
	TanggalKejadian     time.Time         `json:"tanggal_kejadian"`
	KategoriLokasiKasus string            `json:"kategori_lokasi_kasus"`
	AlamatTKP           string            `json:"alamat_tkp"`
	AlamatDetailTKP     string            `json:"alamat_detail_tkp"`
	KronologisKasus     string            `json:"kronologis_kasus"`
	Status              string            `json:"status"`
	AlasanDibatalkan    string            `json:"alasan_dibatalkan"`
	WaktuDilihat        *time.Time        `json:"waktu_dilihat"`
	UserIDMelihat       *uint             `json:"userid_melihat,omitempty"`
	WaktuDiproses       *time.Time        `json:"waktu_diproses"`
	WaktuDibatalkan     *time.Time        `json:"waktu_dibatalkan"`
	Dokumentasi         datatypes.JSONMap `json:"dokumentasi" form:"image" gorm:"type:json"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
}

// package models

// import (
// 	"encoding/json"
// 	"strconv"
// 	"time"

// 	"gorm.io/datatypes"
// )

// type Laporan struct {
// 	NoRegistrasi        string            `gorm:"primaryKey" json:"no_registrasi"`
// 	User                User              `gorm:"foreignKey:UserID"`
// 	UserID              uint              `json:"user_id"`
// 	ViolenceCategory    ViolenceCategory  `gorm:"foreignKey:KategoriKekerasanID"`
// 	KategoriKekerasanID uint              `json:"kategori_kekerasan_id"`
// 	TanggalPelaporan    time.Time         `json:"tanggal_pelaporan"`
// 	TanggalKejadian     time.Time         `json:"tanggal_kejadian"`
// 	KategoriLokasiKasus string            `json:"kategori_lokasi_kasus"`
// 	AlamatTKP           string            `json:"alamat_tkp"`
// 	AlamatDetailTKP     string            `json:"alamat_detail_tkp"`
// 	KronologisKasus     string            `json:"kronologis_kasus"`
// 	Status              string            `json:"status"`
// 	AlasanDibatalkan    string            `json:"alasan_dibatalkan"`
// 	WaktuDilihat        *time.Time        `json:"waktu_dilihat,omitempty"`
// 	UserIDMelihat       *uint             `json:"userid_melihat,omitempty"`
// 	WaktuDiproses       *time.Time        `json:"waktu_diproses,omitempty"`
// 	WaktuDibatalkan     *time.Time        `json:"waktu_dibatalkan,omitempty"`
// 	Dokumentasi         datatypes.JSONMap `json:"dokumentasi" form:"image" gorm:"type:json"`
// 	CreatedAt           time.Time         `json:"created_at"`
// 	UpdatedAt           time.Time         `json:"updated_at"`
// }

// func (l Laporan) MarshalJSON() ([]byte, error) {
// 	type Alias Laporan
// 	aux := &struct {
// 		WaktuDilihat    string `json:"waktu_dilihat"`
// 		UserIDMelihat   string `json:"userid_melihat"`
// 		WaktuDiproses   string `json:"waktu_diproses"`
// 		WaktuDibatalkan string `json:"waktu_dibatalkan"`
// 		*Alias
// 	}{
// 		Alias: (*Alias)(&l),
// 	}

// 	if l.WaktuDilihat != nil {
// 		aux.WaktuDilihat = l.WaktuDilihat.Format(time.RFC3339)
// 	} else {
// 		aux.WaktuDilihat = ""
// 	}

// 	if l.UserIDMelihat != nil {
// 		aux.UserIDMelihat = strconv.FormatUint(uint64(*l.UserIDMelihat), 10)
// 	} else {
// 		aux.UserIDMelihat = ""
// 	}

// 	if l.WaktuDiproses != nil {
// 		aux.WaktuDiproses = l.WaktuDiproses.Format(time.RFC3339)
// 	} else {
// 		aux.WaktuDiproses = ""
// 	}

// 	if l.WaktuDibatalkan != nil {
// 		aux.WaktuDibatalkan = l.WaktuDibatalkan.Format(time.RFC3339)
// 	} else {
// 		aux.WaktuDibatalkan = ""
// 	}

// 	return json.Marshal(aux)
// }
