package models

import (
	"time"

	"gorm.io/gorm"
)

type Vehicle struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Make         string         `json:"make" gorm:"not null"`
	LicensePlate string         `json:"license" gorm:"not null;unique"`
	Model        string         `json:"model" gorm:"not null"`
	Year         int            `json:"year" gorm:"not null"`
	Color        string         `json:"color" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Images       []Image        `json:"images" gorm:"foreignKey:VehicleID"`
}

type Image struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	VehicleID   uint           `json:"vehicle_id" gorm:"not null"`
	Filename    string         `json:"filename" gorm:"not null"`
	Path        string         `json:"path" gorm:"not null"`
	ContentType string         `json:"content_type" gorm:"not null"`
	Size        int64          `json:"size" gorm:"not null"`
	Processed   bool           `json:"processed" gorm:"not null;default:false"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Analayses   []Analysis     `json:"analyses" gorm:"foreignKey:ImageID"`
}

type Analysis struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	ImageID         uint           `json:"image_id" gorm:"not null"`
	VehicleType     string         `json:"vehicle_type" gorm:"not null"`
	ConfidenceScore float64        `json:"confidence" gorm:"not null"`
	LicensePlate    string         `json:"license_plate" gorm:"not null"`
	DamageDetected  bool           `json:"damage_detected" gorm:"not null;default:false"`
	DamageDetails   string         `json:"damage_details" gorm:"not null"`
	ProcessedAt     time.Time      `json:"processed_at" gorm:"autoCreateTime"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
