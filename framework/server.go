package framework

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/img21326/andrew-framework/helper"
	"github.com/img21326/andrew-framework/middleware"
	"github.com/spf13/viper"
)

func ReadConf() {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found, reading from the environment")
		} else {
			fmt.Printf("Fatal error reading config file: %s", err)
		}
	}
}

func InitDB() {
	dbOption := helper.DBOption{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		DBName:   viper.GetString("DB_NAME"),
		SSLMode:  viper.GetString("DB_SSL_MODE"),
	}
	helper.InitDB(dbOption)
}

func InitServer() {
	ReadConf()
	InitDB()
}

func InitGin() *gin.Engine {
	r := gin.Default()

	r.HTMLRender = loadTemplates("templates")

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
		Handler: InitGin(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	workerContext, workerCancel := context.WithCancel(context.Background())
	workerCount := viper.GetInt("WORKER_COUNT")
	workerDoneChans := make([]chan struct{}, workerCount)
	for i := 0; i < workerCount; i++ {
		workerDoneChans[i] = make(chan struct{})
		go helper.StartWorker(workerContext, workerDoneChans[i])
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	workerCancel()
	log.Println("Shutdown Server ...")

	defer cancel()
	defer helper.WaitForLoggerComplete()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	for i := 0; i < workerCount; i++ {
		<-workerDoneChans[i]
	}
	log.Println("Server exiting")
}

func StartAWSLambda() {
	var ginLambda *ginAdapter.GinLambda

	handler := func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return ginLambda.ProxyWithContext(ctx, req)
	}

	r := InitGin()
	ginLambda = ginAdapter.New(r)
	lambda.Start(handler)
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.html")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/includes/*.html")
	if err != nil {
		panic(err.Error())
	}

	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
