package api

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/jarpis-nl/go-base-backend/internal/api/v1"
	"gorm.io/gorm"
)

func NewAPI(engine *gin.Engine, database *gorm.DB) {
	api := engine.Group("/api")
	v1Router := v1.New(database)
	v1Router.Setup(api)
}
