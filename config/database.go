package config

import (
	"fmt"
	"os"

	"github.com/athunlal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBconnect() *gorm.DB {
	dns := os.Getenv("DB_URL")
	DB, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Admin{})
	return DB

}
