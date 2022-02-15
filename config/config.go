package config

import (
	"log"

	"github.com/spf13/viper"
)

type CANType string

const (
	SOCKETCAN CANType = "SOCKETCAN"
	USB       CANType = "USB"
	TCP       CANType = "TCP"
	MOCK      CANType = "MOCK"
)

type Config struct {
	SQLite struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"sqlite"`
	DefaultUser struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"default_user"`
	Lockout struct {
		Enabled    bool
		ListenAddr string `mapstructure:"listen_addr"`
	}
	TimeSeries struct {
		Directory string
	} `mapstructure:"time_series"`
	HTTP struct {
		ListenAddr string `mapstructure:"listen_addr"`
	} `mapstructure:"http"`
	CAN struct {
		CANType       CANType `mapstructure:"can_type"`
		InterfaceName string  `mapstructure:"interface_name"`
		Host          string  `mapstructure:"host"`
		Port          int     `mapstructure:"port"`
		SerialPort    string  `mapstructure:"serial_port"`
		BaudRate      int     `mapstructure:"baud_rate"`
	} `mapstructure:"can"`
}

var CurrentConfig Config

func Init() {
	v := viper.NewWithOptions(viper.KeyDelimiter(":"))
	v.SetDefault("sqlite", map[string]interface{}{
		"dsn": "./main.db",
	})
	v.SetDefault("default_user", map[string]interface{}{
		"username": "test",
		"password": "test",
	})
	v.SetDefault("lockout", map[string]interface{}{
		"enabled":     true,
		"listen_addr": ":5445",
	})
	v.SetDefault("time_series", map[string]interface{}{
		"directory": "./sensor_data",
	})
	v.SetDefault("http", map[string]interface{}{
		"listen_addr": ":8080",
	})
	v.SetDefault("can", map[string]interface{}{
		"can_type":       "MOCK",
		"interface_name": "can0",
		"host":           "localhost",
		"port":           8881,
		"serial_port":    "/dev/tty.usbserial-14140",
		"baud_rate":      500000000,
	})
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		log.Printf("Couldn't read config.yaml, using default settings: %s", err)
	}
	v.Unmarshal(&CurrentConfig)
}
