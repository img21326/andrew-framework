package andrewframework

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/img21326/andrew_framework/helper"
	"github.com/img21326/andrew_framework/middleware"
	"github.com/spf13/viper"
)

func ReadConf() {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s", err)
		return
	}
}

func InitDB() {
	dbOption := helper.DBOption{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
	}
	fmt.Printf("%+v\n", viper.AllSettings())
	helper.InitDB(dbOption)
}

func InitServer() *gin.Engine {
	ReadConf()
	InitDB()
	r := gin.Default()

	r.Use(middleware.WithLoggerMiddleware())
	r.Use(middleware.WithRecoverMiddleware())
	r.Use(middleware.WithRequestLogMiddleware())
	r.Use(middleware.WithGormMiddleware())
	r.Use(middleware.ReturnErrorMiddleware())

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
		})
	})

	for _, router := range RouterList {
		router.AddRoute(r)
	}
	return r
}

func Start() {
	srv := &http.Server{
		Addr:    ":8000",
		Handler: InitServer(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer helper.WaitForLoggerComplete()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

func StartAWSLambda() {
	var ginLambda *ginAdapter.GinLambda

	handler := func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return ginLambda.ProxyWithContext(ctx, req)
	}

	r := InitServer()
	ginLambda = ginAdapter.New(r)
	lambda.Start(handler)
}
