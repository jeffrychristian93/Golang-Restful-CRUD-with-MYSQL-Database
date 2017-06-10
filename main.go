package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/db-example")

type User struct {
	Id int `db:"id"`
	Username 	string `db:"username" form:"username"`
	Password 	string `db:"password" form:"password"`
	FirstName 	string `db:"first_name" form:"first_name"`
	MiddleName 	string `db:"middle_name" form:"middle_name"`
	LastName 	string `db:"last_name" form:"last_name"`
	Email 		string `db:"email" form:"email"`
	MobilePhone string `db:"mobile_phone" form:"mobile_phone"`
	LoginAttempt 	int `db:"login_attempt" form:"login_attempt"`
	RemoteAddress 	string `db:"remote_address" form:"remote_address"`
	ActiveStatus 	int `db:"active_status" form:"active_status"`
}

func getAll(c *gin.Context){
	var (
		user User
		users []User
	)
	rows, err := db.Query("select id, username, first_name, middle_name, last_name, email, mobile_phone, login_attempt, active_status from user;")

	if err != nil {
		fmt.Print(err.Error())
	}

	for rows.Next() {
		rows.Scan(&user.Id, &user.Username, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.MobilePhone, &user.LoginAttempt, &user.ActiveStatus)
		users = append(users, user)
	}

	defer rows.Close()
	c.JSON(http.StatusOK, users)
}

func add(c *gin.Context) {
	Username := c.PostForm("username")
	Password := c.PostForm("password")
	FirstName := c.PostForm("first_name")
	MiddleName := c.PostForm("middle_name")
	LastName := c.PostForm("last_name")
	Email := c.PostForm("email")
	MobilePhone := c.PostForm("mobile_phone")
	LoginAttempt := 0
	RemoteAddress := c.PostForm("remote_address")
	ActiveStatus := 1
	
	if Username == "" || Password == "" || FirstName == "" || Email == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Please fill all mandatory field"),
		})
		return
	}
	
	if isEmailOrUsernameAlreadyExist(Username, Email){
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Username or email is already in use"),
		})
		return
	}
	
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
        panic(err)
    }
	
	stmt, err := db.Prepare("insert into user (username, password, first_name, middle_name, last_name, email, mobile_phone, login_attempt, remote_address, active_status) values(?,?,?,?,?,?,?,?,?,?);")
	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = stmt.Exec(Username, string(HashedPassword), FirstName, MiddleName, LastName, Email, MobilePhone, LoginAttempt, RemoteAddress, ActiveStatus)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer stmt.Close()

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("successfully created"),
	})
}

func isEmailOrUsernameAlreadyExist(username, email string) bool {
	var Id int
	err := db.QueryRow("SELECT id FROM user WHERE username= ? OR email= ?", username, email).Scan(&Id)
	
	switch {
	case err == sql.ErrNoRows:
			log.Printf("No user with that ID.")
			return false
	case err != nil:
			log.Fatal(err)
			return false
	default:
		log.Printf("Found")
		return true
	}
}

func isPasswordMatch(id int, pwd string) bool {
	var passwordFromDB string
	
	row := db.QueryRow("select password from user where id = ?;", id)	
	err = row.Scan(&passwordFromDB)
	
	passwordInByte  := []byte(passwordFromDB)
	old_password := []byte(pwd)
	
	err = bcrypt.CompareHashAndPassword(passwordInByte, old_password)
	if err == nil {
		//If nil => "password match"
		return true
	} else {
		return false
	}
}

func update(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	Username := c.PostForm("username")
	Password := c.PostForm("password")
	OldPassword := c.PostForm("old_password")
	FirstName := c.PostForm("first_name")
	MiddleName := c.PostForm("middle_name")
	LastName := c.PostForm("last_name")
	Email := c.PostForm("email")
	MobilePhone := c.PostForm("mobile_phone")
	LoginAttempt := 0
	RemoteAddress := c.PostForm("remote_address")
	ActiveStatus := 1
	
	if Username == "" || Password == "" || OldPassword == "" || FirstName == "" || Email == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Please fill all mandatory field"),
		})
		return
	}
	
	if !isPasswordMatch(Id, OldPassword) {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Old password does not match"),
		})
		return
	}

	var HashedPassword []byte
	
	HashedPassword, err = bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
        panic(err)
    }
	
	stmt, err := db.Prepare("update user set username= ?, password= ?, first_name= ?, middle_name= ?, last_name= ?, email= ?, mobile_phone= ?, login_attempt= ?, remote_address= ?, active_status= ? where id= ?;")
	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = stmt.Exec(Username, string(HashedPassword), FirstName, MiddleName, LastName, Email, MobilePhone, LoginAttempt, RemoteAddress, ActiveStatus, Id)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer stmt.Close()

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully updated"),
	})
}

func getById(c *gin.Context) {
	var user User 

	id := c.Param("id")
	row := db.QueryRow("select id, username, first_name, middle_name, last_name, email, mobile_phone, login_attempt, active_status from user where id = ?;", id)
	
	err = row.Scan(&user.Id, &user.Username, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.MobilePhone, &user.LoginAttempt, &user.ActiveStatus)
	if err != nil {
		c.JSON(http.StatusOK, nil)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func delete(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("delete from user where id= ?;")
	
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec(id)
	
	if err != nil {
		fmt.Print(err.Error())
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully deleted user with ID : %s", id),
	})
}

func createTable(){
	stmt, err := db.Prepare("CREATE TABLE user (id int NOT NULL AUTO_INCREMENT, username varchar(40), password varchar(255), first_name varchar(40), middle_name varchar(40), last_name varchar(40), email varchar(60), mobile_phone varchar(15), login_attempt int(1), remote_address varchar(40), active_status int(1), PRIMARY KEY (id));")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Table is successfully created....")
	}
}

func main() {
	createTable();
	
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}
	
	router := gin.Default()
	router.GET("/api/user/:id", getById)
	router.GET("/api/users", getAll)
	router.POST("/api/user", add)
	router.PUT("/api/user/:id", update)
	router.DELETE("/api/user/:id", delete)
	router.Run(":8000")
}