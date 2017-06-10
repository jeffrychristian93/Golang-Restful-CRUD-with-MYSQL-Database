# Golang-Restful-CRUD-with-MYSQL-Database

![Go Language](http://marcio.io/img/gopher.png?raw=true)

Sample Go Language for CRUD with restful service API

Installation

- Download main.go
- Edit your database connection on this code:
  - sql.Open("mysql", "root@tcp(127.0.0.1:3306)/db-example")
  - sql.Open("mysql", "username:password@tcp(yourIpAddress:port)/yourDataabse")
- Run using this command -> go run main.go
- Or use command -> "go build main.go" to generate executable file

______________________________________________________________________________________

If you got some error, please download the libraries first.

- go get github.com/gin-gonic/gin" -> Gin is a web framework written in Go (Golang)
- go get github.com/go-sql-driver/mysql" -> Mysql driver for Go
- Special for Bcrypt use : git clone https://go.googlesource.com/crypto 

______________________________________________________________________________________

Listening and serving HTTP on :8000
- GET    /api/user/:id             --> main.getById
- GET    /api/users                --> main.getAll
- POST   /api/user                 --> main.add
- PUT    /api/user/:id             --> main.update
- DELETE /api/user/:id             --> main.delete

______________________________________________________________________________________

Source :
- [Go Languange](https://golang.org/)
- [Gin Framework](https://gin-gonic.github.io/gin/)
- [Mysql Database](https://www.mysql.com/)

Other projects in GO :
- [Shortening URL](https://github.com/jeffrychristian93/Golang-Shortening-URL-Mysql)

Thank you..~
______________________________________________________________________________________
