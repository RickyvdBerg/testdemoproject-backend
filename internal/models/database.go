package models

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/jarpis-nl/go-base-backend/internal/utils"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" swaggerignore:"true" json:"-"`
}


type User struct {
	Model
	Email    string `gorm:"uniqueIndex" json:"email"`
	Password string `json:"-" swaggerignore:"true"`
	Name     string `json:"name"`
}

type UserClaims struct {
	*jwt.StandardClaims
	Data User `json:"user"`
}

func (u User) GenerateJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(8 * time.Hour).Unix(),
		},
		Data: u,
	})

	return token.SignedString([]byte(viper.GetString("JWT_TOKEN")))
}

func (u *User) Update(updates UserUpdate) {
	if len(updates.Name) > 0 {
		u.Name = updates.Name
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate password for user
	if len(u.Password) > 0 {
		pass, err := utils.HashPassword(u.Password)

		if err != nil {
			return err
		}

		tx.Statement.SetColumn("Password", pass)
	}

	return nil
}

//This is a collection of database models that need to be automatically migrated

var DatabaseModels = []interface{}{ &User{} }
