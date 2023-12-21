package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"log"
	"os"
	"reflect"
)

var Instance *AppConfig

type AppConfig struct {
	Database Database
	Security Security
}

type Database struct {
	Host     string `mapstructure:"APP_DB_HOST" validate:"required"`
	Port     int    `mapstructure:"APP_DB_PORT" validate:"required"`
	Name     string `mapstructure:"APP_DB_DB" validate:"required"`
	Username string `mapstructure:"APP_DB_USERNAME" validate:"required"`
	Password string `mapstructure:"APP_DB_PASSWORD" validate:"required"`
}

type Security struct {
	JwtSecKey string `mapstructure:"APP_SEC_JWTKEY" validate:"required"`
	JwtTTL    string `mapstructure:"APP_SEC_JWT_TTL" validate:"required"`
}

func InitConfigInstance() (err error) {
	dbCfg, _ := scanConfig[Database]()
	secCfg, _ := scanConfig[Security]()

	Instance = &AppConfig{
		Database: dbCfg,
		Security: secCfg,
	}
	log.Printf("Config %+v", Instance)

	return
}

func getOsEnv(name string) (key string, value interface{}) {
	key = name
	value = os.Getenv(name)
	return
}

func scanConfig[T any]() (t T, err error) {
	tF := reflect.TypeOf(t)
	for i := 0; i < tF.NumField(); i++ {
		f := tF.Field(i)
		tag := f.Tag.Get("mapstructure")
		viper.Set(getOsEnv(tag))
	}

	err = viper.Unmarshal(&t)
	if err != nil {
		panic(fmt.Errorf("err failed to unmarshal config [%+v]: %+v", tF, err))
	}

	validate := validator.New()
	err = validate.Struct(t)
	if err != nil {
		panic(fmt.Errorf("err invalid configuration %s", err))
	}

	return
}
