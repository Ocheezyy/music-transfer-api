package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Ocheezyy/music-transfer-api/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (ac *AuthController) Login(c *gin.Context) {

	var authInput models.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	ac.DB.Where("email=?", authInput.Email).Find(&userFound)

	if userFound.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userFound.ID,
		"email": userFound.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}

func (ac *AuthController) CreateUser(c *gin.Context) {

	var authInput models.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		log.Printf("CreateUser bad request: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	ac.DB.Where("email=?", authInput.Email).Find(&userFound)

	if userFound.ID != 0 {
		log.Print("CreateUser error: already exists")
		c.JSON(http.StatusConflict, gin.H{"error": "email already used"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("CreateUser failed password hash: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:    authInput.Email,
		Password: string(passwordHash),
	}

	ac.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (ac *AuthController) GetUserProfile(c *gin.Context) {
	user, _ := c.Get("currentUser")

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
