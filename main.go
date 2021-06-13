package main

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/jarpis-nl/go-base-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-contrib/cors"

	"github.com/jarpis-nl/go-base-backend/internal/api"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	_ "github.com/jarpis-nl/go-base-backend/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth
func initConfig() {
	// Set config type of config files
	viper.SetConfigType("yaml")

	// Set config path
	viper.AddConfigPath(".")

	// Set name of config file (without extension)
	viper.SetConfigName("config")

	// Automatically pick up env variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			logrus.Fatal("Unable to find config file, make sure you have provided your config file based on the config.base.yaml in the root of the project")
		}

		// Something else went wrong
		logrus.Fatalf("Something went wrong: %s", err)
	}
}

func initLogger() *logrus.Logger {
	logger := logrus.New()

	logger.SetOutput(os.Stdout)

	return logger
}

func setupGin(logger *logrus.Logger) *gin.Engine {
	// Default GIN configuration
	router := gin.New()

	//router.Use(ginlogrus.Logger(logger))

	router.Use(gin.Logger())

	router.Use(gin.Recovery())

	// Enable CORS if specified
	if viper.GetBool("SERVER_CORS") {
		router.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"*"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	// Create a static folder that allows users to read files from the server
	router.Static("/assets", "./assets")

	return router
}

func setupDB() *gorm.DB {
	sslmode := "disable"

	if viper.GetBool("DATABASE_SSL") {
		sslmode = "enable"
	}

	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		viper.GetString("DATABASE_USER"),
		url.QueryEscape(viper.GetString("DATABASE_PASSWORD")),
		viper.GetString("DATABASE_HOST"),
		viper.GetString("DATABASE_DATABASE"),
		viper.GetInt("DATABASE_PORT"),
		sslmode,
		viper.GetString("DATABASE_TIMEZONE"),
	)

	if viper.IsSet("DATABASE_URL") {
		dsn = viper.GetString("DATABASE_URL")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logrus.Panic(err)
	}

	if err := db.AutoMigrate(models.DatabaseModels...); err != nil {
		logrus.Panic(err)
	}

	return db
}

func main() {
	logger := initLogger()
	initConfig()
	db := setupDB()

	r := setupGin(logger)
	
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	api.NewAPI(r, db)

	// Create host address for GIN to run on using config
	routerAddress := fmt.Sprintf("%s:%d", viper.GetString("SERVER_ADDRESS"), viper.GetInt("SERVER_PORT"))

	r.Run(routerAddress)
}
