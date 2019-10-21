package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

type User struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Height float64 `json:"height"`
}

var db *sql.DB

const (
	dbhost = "DBHOST"
	dbport = "DBPORT"
	dbuser = "DBUSER"
	dbpass = "DBPASS"
	dbname = "DBNAME"
)

func initDb() {
	config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbpass], config[dbname])

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
}

func dbConfig() map[string]string {
	conf := make(map[string]string)
	host, ok := os.LookupEnv(dbhost)
	if !ok {
		panic("DBHOST environment variable required but not set")
	}
	port, ok := os.LookupEnv(dbport)
	if !ok {
		panic("DBPORT environment variable required but not set")
	}
	user, ok := os.LookupEnv(dbuser)
	if !ok {
		panic("DBUSER environment variable required but not set")
	}
	password, ok := os.LookupEnv(dbpass)
	if !ok {
		panic("DBPASS environment variable required but not set")
	}
	name, ok := os.LookupEnv(dbname)
	if !ok {
		panic("DBNAME environment variable required but not set")
	}
	conf[dbhost] = host
	conf[dbport] = port
	conf[dbuser] = user
	conf[dbpass] = password
	conf[dbname] = name
	return conf
}

func deleteUser(c echo.Context) error {
	// user := User{}
	sqlStatement := `DELETE FROM users WHERE id = $1;`
	rows, err := db.Query(sqlStatement, c.Param("id"))
	if err != nil {
		return err
	}
	defer rows.Close()
	// rows.Next()
	return c.String(http.StatusOK, "WORKED!")
}

func getUsers(c echo.Context) error {
	result := ""
	sqlStatement := `SELECT * FROM users;`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Height)
		if err != nil {
			panic(err)
		}
		result = result + fmt.Sprintf("ID: %v, Name: %v, Height: %v, Age: %v\n", user.Id, user.Name, user.Height, user.Age)
	}
	return c.String(http.StatusOK, result)
}

func createUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	id := u.Id
	name := u.Name
	age := u.Age
	height := u.Height
	sqlStatement := `INSERT INTO users VALUES ($1, $2, $3, $4);`
	_, err := db.Exec(sqlStatement, id, name, age, height)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, u)
}
func getUser(c echo.Context) error {
	result := ""
	sqlStatement := `SELECT * FROM users WHERE id = $1;`
	rows, err := db.Query(sqlStatement, c.Param("id"))
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Height)
		if err != nil {
			panic(err)
		}
		result = result + fmt.Sprintf("ID: %v, Name: %v, Height: %v, Age: %v\n", user.Id, user.Name, user.Height, user.Age)
	}
	return c.String(http.StatusOK, result)
}
func updateUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	sqlStatement := `UPDATE users SET name=$1, age=$2, height=$3 WHERE id=$4;`
	rows, err := db.Query(sqlStatement, u.Name, u.Age, u.Height, c.Param("id"))
	if err != nil {
		return err
	}
	defer rows.Close()
	return c.String(http.StatusOK, "WORKED!")
}

func main() {
	os.Setenv(dbhost, "127.0.0.1")
	os.Setenv(dbport, "5432")
	os.Setenv(dbuser, "joao")
	os.Setenv(dbpass, "123456789")
	os.Setenv(dbname, "users")

	initDb()
	defer db.Close()

	e := echo.New()
	e.GET("/api/users/:id", getUser)
	e.GET("/api/users", getUsers)
	e.POST("/api/users", createUser)
	e.DELETE("/api/users/:id", deleteUser)
	e.PUT("/api/users/:id", updateUser)
	log.Fatal(http.ListenAndServe("localhost:8080", e))
}


