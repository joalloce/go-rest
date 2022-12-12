package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Hello,World")
	r := setRouter()
	r.Run()
}

func setRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pang",
		})
	})

	userRepo := New()
	r.POST("/users", userRepo.CreateUser)
	r.GET("/users", userRepo.GetUsers)
	r.GET("/users/:id", userRepo.GetUser)
	r.PUT("/users/:id", userRepo.UpdateUser)
	r.DELETE("/users/:id", userRepo.DeleteUser)

	return r
}

const DB_USERNAME = "root"
const DB_PASSWORD = "password"
const DB_NAME = "db"
const DB_HOST = "db"
const DB_PORT = "3306"

var Db *gorm.DB

func initDb() *gorm.DB {
	Db = connectDB()
	return Db
}

func connectDB() *gorm.DB {
	var err error
	dsn := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?" + "parseTime=true&loc=Local"
	fmt.Println("dsn : ", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Error connecting to database :", err)
		return nil
	}

	return db
}

type User struct {
	gorm.Model
	ID        int
	FirstName string
	LastName  string
	Email     string
}

// create a user
func CreateUserDB(db *gorm.DB, User *User) (err error) {
	err = db.Create(User).Error
	if err != nil {
		return err
	}
	return nil
}

// get users
func GetUsersDB(db *gorm.DB, User *[]User) (err error) {
	err = db.Find(User).Error
	if err != nil {
		return err
	}
	return nil
}

// get user by id
func GetUserDB(db *gorm.DB, User *User, id int) (err error) {
	err = db.Where("id = ?", id).First(User).Error
	if err != nil {
		return err
	}
	return nil
}

// update user
func UpdateUserDB(db *gorm.DB, User *User) (err error) {
	db.Save(User)
	return nil
}

// delete user
func DeleteUserDB(db *gorm.DB, User *User, id int) (err error) {
	db.Where("id = ?", id).Delete(User)
	return nil
}

type UserRepo struct {
	Db *gorm.DB
}

func New() *UserRepo {
	db := initDb()
	db.AutoMigrate(&User{})
	return &UserRepo{Db: db}
}

// create user
func (repository *UserRepo) CreateUser(c *gin.Context) {
	var user User
	c.BindJSON(&user)
	err := CreateUserDB(repository.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}

// get users
func (repository *UserRepo) GetUsers(c *gin.Context) {
	var user []User
	err := GetUsersDB(repository.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}

// get user by id
func (repository *UserRepo) GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user User
	err := GetUserDB(repository.Db, &user, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}

// update user
func (repository *UserRepo) UpdateUser(c *gin.Context) {
	var user User
	id, _ := strconv.Atoi(c.Param("id"))
	err := GetUserDB(repository.Db, &user, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.BindJSON(&user)
	err = UpdateUserDB(repository.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}

// delete user
func (repository *UserRepo) DeleteUser(c *gin.Context) {
	var user User
	id, _ := strconv.Atoi(c.Param("id"))
	err := DeleteUserDB(repository.Db, &user, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
