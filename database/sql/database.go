package sql

import (
	"log"

	"github.com/Liquid-Propulsion/mainland-server/config"
	"github.com/Liquid-Propulsion/mainland-server/types"
	"github.com/matthewhartstonge/argon2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Init(dsn string) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
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
	Database.AutoMigrate(&types.IslandNode{})

	argon := argon2.DefaultConfig()
	hash, err := argon.HashEncoded([]byte(config.CurrentConfig.DefaultUser.Password))
	if err == nil {
		userType := types.User{
			Name:         config.CurrentConfig.DefaultUser.Username,
			Username:     config.CurrentConfig.DefaultUser.Username,
			PasswordHash: hash,
			TOTPEnabled:  false,
		}
		Database.Create(&userType)
	}
}
