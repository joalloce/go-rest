package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	r := setRouter()
	r.Run()
}

func setRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	userRepo := New()

	r.POST("/users", userRepo.createUser)
	r.GET("/users", userRepo.getUsers)
	r.GET("/users/:id", userRepo.getUser)
	r.PUT("/users/:id", userRepo.updateUser)
	r.DELETE("/users/:id", userRepo.deleteUser)

	return r
}

const DB_USERNAME = "root"
const DB_PASSWORD = "password"
const DB_NAME = "db"
const DB_HOST = "db"
const DB_PORT = "3306"

// db connection
func connectDB() *gorm.DB {
	var err error
	DB_URI := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?" + "parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(DB_URI), &gorm.Config{})

	if err != nil {
		fmt.Println("Error db :", err)
		return nil
	}

	return db
}

// User model
type User struct {
	gorm.Model
	ID        int
	FirstName string
	LastName  string
	Email     string
}

type UserRepo struct {
	Db *gorm.DB
}

func New() *UserRepo {
	db := connectDB()
	db.AutoMigrate(&User{})
	return &UserRepo{Db: db}
}

// create user
func (repository *UserRepo) createUser(context *gin.Context) {
	var user User
	err := context.BindJSON(&user) // take the data from body
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	err = repository.Db.Create(&user).Error // create user
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, user)
}

// get users
func (repository *UserRepo) getUsers(context *gin.Context) {
	var user []User
	err := repository.Db.Find(&user).Error // get all users
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, user)
}

// get user by id
func (repository *UserRepo) getUser(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id")) // take id from params
	var user User
	err := repository.Db.Where("id = ?", id).First(&user).Error // find user
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, user)
}

// update user
func (repository *UserRepo) updateUser(context *gin.Context) {
	var user User
	id, _ := strconv.Atoi(context.Param("id"))                  // take id from params
	err := repository.Db.Where("id = ?", id).First(&user).Error // find user
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	err = context.BindJSON(&user) // take the data from body
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	err = repository.Db.Save(&user).Error // update user
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, user)
}

// delete user
func (repository *UserRepo) deleteUser(context *gin.Context) {
	var user User
	id, _ := strconv.Atoi(context.Param("id"))                   // take id from params
	err := repository.Db.Where("id = ?", id).Delete(&user).Error // delete user
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
