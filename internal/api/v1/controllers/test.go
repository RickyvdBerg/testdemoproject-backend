package controllers

import (
	"net/http"

	"github.com/jarpis-nl/go-base-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jarpis-nl/go-base-backend/internal/models"
	"gorm.io/gorm"
)

type TestController struct {
	db *gorm.DB
}

func NewTestController() *TestController {
	return &TestController{}
}

// WhoAmI godoc
// @Summary Based on the authentication token, receive your current information
// @Success 200 {object} models.User
// @Accept json
// @Produce json
// @tags User
// @Security BasicAuth
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Router /whoami [get]
func (t TestController) TestConnection(c *gin.Context) {

	var testResponse models.TestResponse
	testResponse.Message = "Api works"

	utils.SetAPISuccess(c, http.StatusOK, testResponse)
}
