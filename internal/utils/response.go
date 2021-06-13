package utils

import (
	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

type PaginatedResponse struct {
	Data   interface{} `json:"data"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	Total  int64       `json:"total"`
}

type APIErrorResponse struct {
	Error HTTPError `json:"error"`
}

type APISuccessResponse struct {
	Data interface{} `json:"data"`
}

func SetAPIError(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, APIErrorResponse{
		Error: HTTPError{
			Code:    code,
			Message: err.Error(),
		},
	})
}

func SetAPISuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, APISuccessResponse{Data: data})
}

func SetAPIPaginatedSuccess(c *gin.Context, code int, response PaginatedResponse) {
	c.JSON(code, response)
}
