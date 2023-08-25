package helper

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormOption struct {
	Ctx *gin.Context
}

var DB *sql.DB

var once sync.Once

type DBOption struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func InitDB(option DBOption) {
	once.Do(func() {
		var err error
		DB, err = sql.Open("pgx", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Taipei", option.Host, option.Port, option.User, option.Password, option.DBName))
		if err != nil {
			panic(err)
		}
		err = DB.Ping()
		if err != nil {
			panic(err)
		}
		DB.SetConnMaxIdleTime(1 * time.Minute)
		DB.SetConnMaxLifetime(5 * time.Minute)
		DB.SetMaxIdleConns(3)
		DB.SetMaxOpenConns(30)
	})
}

func NewGorm(option GormOption) *gorm.DB {
	config := &gorm.Config{}
	logger := GetLogger(option.Ctx)
	config.Logger = logger
	gorm, err := gorm.Open(postgres.New(postgres.Config{Conn: DB}), config)
	if err != nil {
		panic(err)
	}
	for _, f := range gormHook {
		f(gorm, option.Ctx)
	}
	return gorm
}

func GetGorm(ctx *gin.Context) *gorm.DB {
	return ctx.MustGet("gormDB").(*gorm.DB)
}
