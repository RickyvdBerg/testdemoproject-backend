package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jarpis-nl/go-base-backend/internal/api/v1/controllers"
	"gorm.io/gorm"

  
  "github.com/dgrijalva/jwt-go"
	"github.com/jarpis-nl/go-base-backend/internal/models"
	"github.com/jarpis-nl/go-base-backend/internal/utils"
	"github.com/spf13/viper"
  
)

type Router struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Router {
	return &Router{
		db: db,
	}
}

func (r Router) Setup(routerGroup *gin.RouterGroup) {
	v1 := routerGroup.Group("/v1")
	{
    
		auth := v1.Group("/auth")
		{
			authController := controllers.NewAuthController(r.db)

			auth.POST("/login", authController.Login)
			auth.POST("/register", authController.Register)
		}

		users := v1.Group("/users", r.AuthMiddleware())
		{
			userController := controllers.NewUserController(r.db)

			users.GET("", userController.GetAllUsers)
			users.POST("", userController.CreateUser)
			users.GET("/:id", userController.GetUser)
			users.PUT("/:id", userController.UpdateUser)
			users.DELETE("/:id", userController.DeleteUser)
		}
    

		testController := controllers.NewTestController()
		v1.GET("/test", testController.TestConnection)
	}
}

// AuthMiddleware checks the header for Authorization tokens
// If the request comes from swagger directly it will check for the basic authentication and use that for further authentication
// For anything else it expects the Bearer {token} and checks it validity before continuing
// Finally sets USER_ID of the logged in user in the context so it is accessible for all other handlers
func (r Router) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userId uint
		referer := c.Request.Header.Get("Referer")

		// Check if referer is Swagger for additional security on basic authentication
		if strings.Contains(referer, "swagger") {
			email, password, ok := c.Request.BasicAuth()
			var credentials models.LoginRequest

			if ok {
				credentials = models.LoginRequest{
					Email:    email,
					Password: password,
				}

				var user models.User

				result := r.db.First(&user, "email = ?", credentials.Email)

				// User not found
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					utils.SetAPIError(c, http.StatusUnauthorized, errors.New("invalid username/password"))
					return
				}

				if ok := utils.VerifyPassword(user.Password, credentials.Password); !ok {
					utils.SetAPIError(c, http.StatusUnauthorized, errors.New("invalid token"))
					return
				}

				userId = user.ID
			} else {
				utils.SetAPIError(c, http.StatusUnauthorized, errors.New("invalid token"))
				return
			}
		} else {
			// Process bearer JWT
			bearer := c.GetHeader("Authorization")
			parts := strings.Split(bearer, " ")

			if len(parts) != 2 {
				utils.SetAPIError(c, http.StatusUnauthorized, errors.New("invalid token"))
				return
			} else {
				id, err := isTokenValid(parts[1])
				userId = uint(id)

				if err != nil {
					utils.SetAPIError(c, http.StatusUnauthorized, errors.New("invalid token"))
					return
				}
			}
		}

		c.Set("USER_ID", int(userId))
		c.Next()
	}
}

// isTokenValid checks validity of a bearer token
func isTokenValid(token string) (int, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &models.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(viper.GetString("JWT_TOKEN")), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := parsedToken.Claims.(*models.UserClaims); ok && parsedToken.Valid {
		if claims.Data.ID > 0 {
			return int(claims.Data.ID), nil
		}

		return 0, nil
	} else {
		fmt.Println(err)
	}

	return 0, nil
}

