package helper

import (
	"database/sql"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormOption struct {
	Logger *Logger
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
		DB, err = sql.Open("pgx", "host="+option.Host+" port="+option.Port+" user="+option.User+" password="+option.Password+" dbname="+option.DBName+" sslmode=disable TimeZone=Asia/Shanghai")
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
	if option.Logger != nil {
		config.Logger = option.Logger
	}
	gorm, err := gorm.Open(postgres.New(postgres.Config{Conn: DB}), config)
	if err != nil {
		panic(err)
	}
	return gorm
}

func GetGorm(ctx *gin.Context) *gorm.DB {
	return ctx.MustGet("gormDB").(*gorm.DB)
}
