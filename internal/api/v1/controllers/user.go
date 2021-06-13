package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jarpis-nl/go-base-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jarpis-nl/go-base-backend/internal/models"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{db: db}
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get all users as a paginated resource
// @Success 200 {object} utils.PaginatedResponse
// @Accept json
// @Produce json
// @Security BasicAuth
// @tags User
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Router /users [get]
func (u UserController) GetAllUsers(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if err != nil || offset < 0 {
		utils.SetAPIError(c, http.StatusBadRequest, errors.New("invalid offset"))
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "0"))

	if err != nil || limit < 0 {
		utils.SetAPIError(c, http.StatusBadRequest, errors.New("invalid limit"))
		return
	}

	var total int64

	u.db.Model(&models.User{}).Count(&total)

	var users []models.User

	u.db.Limit(limit).Offset(offset).Find(&users)

	utils.SetAPIPaginatedSuccess(c, http.StatusOK, utils.PaginatedResponse{
		Data:   users,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	})
}

// GetUser godoc
// @Summary Get user
// @Description Get a single user using the ID
// @Success 200 {object} models.User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @tags User
// @Security BasicAuth
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Router /users/{id} [get]
func (u UserController) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		utils.SetAPIError(c, http.StatusBadRequest, errors.New("invalid ID"))
		return
	}

	var user models.User

	result := u.db.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		utils.SetAPIError(c, http.StatusNotFound, errors.New("not found"))
		return
	}

	utils.SetAPISuccess(c, http.StatusOK, user)
}

// CreateUser godoc
// @Summary Create user
// @Description Create a new user
// @Success 201 {object} models.User
// @Accept json
// @Produce json
// @Security BasicAuth
// @Param user body models.UserCreate true "User"
// @tags User
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Router /users [post]
func (u UserController) CreateUser(c *gin.Context) {
	var createUser models.UserCreate

	if err := c.ShouldBindJSON(&createUser); err != nil {
		utils.SetAPIError(c, http.StatusBadRequest, err)
		return
	}

	newUser := models.User{
		Email:    createUser.Email,
		Password: createUser.Password,
	}

	dbResult := u.db.Create(&newUser)

	if dbResult.Error != nil {
		utils.SetAPIError(c, http.StatusBadRequest, dbResult.Error)
		return
	}

	utils.SetAPISuccess(c, http.StatusCreated, newUser)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update a user with new values
// @Success 200 {object} models.User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @tags User
// @Security BasicAuth
// @Param user body models.UserUpdate true "User"
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Router /users/{id} [put]
func (u UserController) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		utils.SetAPIError(c, http.StatusBadRequest, errors.New("invalid ID"))
		return
	}

	var user models.UserUpdate

	if err := c.ShouldBindJSON(&user); err != nil {
		utils.SetAPIError(c, http.StatusBadRequest, err)
		return
	}

	var us models.User

	u.db.First(&us, id)

	us.Update(user)

	result := u.db.Save(&us)

	if result.Error != nil {
		utils.SetAPIError(c, http.StatusBadRequest, result.Error)
		return
	}

	utils.SetAPISuccess(c, http.StatusOK, us)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user
// @Success 204
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @tags User
// @Security BasicAuth
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Router /users/{id} [delete]
func (u UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		utils.SetAPIError(c, http.StatusBadRequest, errors.New("invalid ID"))
		return
	}

	if id == c.GetInt("USER_ID") {
		utils.SetAPIError(c, http.StatusBadRequest, errors.New("do not delete yourself"))
		return
	}

	result := u.db.Delete(&models.User{}, id)

	if result.Error != nil {
		utils.SetAPIError(c, http.StatusInternalServerError, errors.New("unable to delete user"))
		return
	}

	c.Status(http.StatusNoContent)
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
func (u UserController) WhoAmI(c *gin.Context) {
	// Set by auth middleware
	userId := c.GetInt("USER_ID")

	var user models.User

	result := u.db.First(&user, userId)

	if result.Error != nil {
		utils.SetAPIError(c, http.StatusBadRequest, errors.New("unable to find user"))
		return
	}

	utils.SetAPISuccess(c, http.StatusOK, user)
}
