package models

import (
	"github.com/golang-jwt/jwt"
	"gorm.io/datatypes"
	"time"
)

type User struct {
	ID           uint              `gorm:"primaryKey" json:"id"`
	FullName     string            `json:"full_name"`
	Username     string            `json:"username" gorm:"size:255;unique;not null"`
	Role         string            `json:"role" gorm:"type:enum('masyarakat','admin');default:'masyarakat'"`
	PhotoProfile string            `json:"photo_profile" gorm:"default:null"`
	PhoneNumber  string            `json:"phone_number" gorm:"unique;not null"`
	Email        string            `json:"email" gorm:"size:255;unique;not null"`
	NIK          uint              `json:"nik"`
	TempatLahir  string            `json:"tempat_lahir" gorm:"default:null"`
	TanggalLahir time.Time         `json:"tanggal_lahir" gorm:"default:null"`
	JenisKelamin string            `json:"jenis_kelamin" gorm:"default:null"`
	Alamat       datatypes.JSONMap `json:"alamat" gorm:"type:json"`
	Password     string            `json:"password"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type LoginCredentials struct {
	Email       string `json:"email"`
	Password    string `json:"password" binding:"required"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
}

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}
