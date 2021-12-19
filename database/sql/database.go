package sql

import (
	"log"

	"github.com/Liquid-Propulsion/mainland-server/config"
	"github.com/Liquid-Propulsion/mainland-server/types"
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
	Database.AutoMigrate(&types.Stage{})
	Database.AutoMigrate(&types.SafetyCheck{})
	Database.AutoMigrate(&types.SOP{})
	Database.AutoMigrate(&types.Solenoid{})
	Database.AutoMigrate(&types.User{})
	Database.AutoMigrate(&types.Sensor{})
}
