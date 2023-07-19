package helper

import (
	"database/sql"
	"os"
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

func init() {
	once.Do(func() {
		var err error
		DB, err = sql.Open("pgx", os.Getenv("DATABASE_URL"))
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
	return ctx.MustGet("gorm").(*gorm.DB)
}
