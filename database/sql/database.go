package sql

import (
	"log"

	"github.com/Liquid-Propulsion/mainland-server/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Init() {
	db, err := gorm.Open(sqlite.Open(config.CurrentConfig.SQLite.DSN), &gorm.Config{})
	if err != nil {
		log.Panicf("Couldn't open the database.")
	}
	Database = db
	Database.AutoMigrate(&Stage{})
}
