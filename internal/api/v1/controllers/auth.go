package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jarpis-nl/go-base-backend/internal/models"
	"github.com/jarpis-nl/go-base-backend/internal/utils"
	"gorm.io/gorm"
)

type AuthController struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db: db}
}

// Login godoc
// @Summary Login
// @Success 200 {object} models.TokenResponse
// @Accept json
// @Produce json
// @tags Authentication
// @Param credentials body models.LoginRequest true "Credentials"
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Router /auth/login [post]
func (ac AuthController) Login(c *gin.Context) {
	var credentials models.LoginRequest

	if err := c.ShouldBindJSON(&credentials); err != nil {
		utils.SetAPIError(c, http.StatusBadRequest, err)
		return
	}

	var user models.User

	result := ac.db.Where("email = ?", credentials.Email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			utils.SetAPIError(c, http.StatusUnauthorized, errors.New("invalid email/password"))
		} else {
			utils.SetAPIError(c, http.StatusUnauthorized, errors.New("an unknown error occurred"))
		}
		return
	}

	if ok := utils.VerifyPassword(user.Password, credentials.Password); !ok {
		utils.SetAPIError(c, http.StatusUnauthorized, errors.New("invalid email/password"))
		return
	}

	tokenString, err := user.GenerateJWT()

	if err != nil {
		utils.SetAPIError(c, http.StatusInternalServerError, errors.New("something went wrong, please try again later"))
		return
	}

	response := models.TokenResponse{Token: tokenString}

	utils.SetAPISuccess(c, http.StatusOK, response)
}

// Register godoc
// @Summary Register
// @Success 201 {object} models.User
// @Accept json
// @Produce json
// @tags Authentication
// @Param user body models.UserCreate true "User"
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Router /auth/register [post]
func (ac AuthController) Register(c *gin.Context) {
	var createUser models.UserCreate

	if err := c.ShouldBindJSON(&createUser); err != nil {
		utils.SetAPIError(c, http.StatusBadRequest, err)
		return
	}

	newUser := models.User{
		Email:    createUser.Email,
		Password: createUser.Password,
	}

	dbResult := ac.db.Create(&newUser)

	if dbResult.Error != nil {
		utils.SetAPIError(c, http.StatusBadRequest, dbResult.Error)
		return
	}

	utils.SetAPISuccess(c, http.StatusCreated, newUser)
}
