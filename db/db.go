package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Humidity struct {
	gorm.Model
	Value  float32 `json:"value"`
	Device string  `json:"device"`
}

type Temperature struct {
	gorm.Model
	Value  float32 `json:"value"`
	Device string  `json:"device"`
}

func NewSqliteCLient() (*gorm.DB, error) {
	dsn := "root:my_password@tcp(127.0.0.1:6603)/mysql?charset=utf8mb4&parseTime=True&loc=Local"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
