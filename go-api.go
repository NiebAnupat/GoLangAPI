package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// album represents data about a record album.
type Attractions struct {
	Id         string `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	Detail     string `db:"detail" json:"detail"`
	Coverimage string `db:"coverimage" json:"coverimage"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/goapi")
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/attractions", getAttractions)
	router.POST("/attractions", addAttractions)
	router.Run("localhost:8080")
	fmt.Println("Server is running on port 8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAttractions(c *gin.Context) {
	var attractions []Attractions

	rows, err := db.Query("select id, name,detail,coverimage from attractions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var a Attractions
		err := rows.Scan(&a.Id, &a.Name, &a.Detail, &a.Coverimage)
		if err != nil {
			log.Fatal(err)
		}
		attractions = append(attractions, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, attractions)
}

// addAttractions adds an album from JSON received in the request body.
func addAttractions(c *gin.Context) {
	var a Attractions

	if err := c.ShouldBind(&a); err != nil {
		log.Fatal(err)
		return
	}

	// log.Fatal(a)

	stmt, err := db.Prepare("insert into attractions (id, name, detail, coverimage) values (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// a.Id = "14"
	// a.Name = "test"
	// a.Detail = "test"
	// a.Coverimage = "test"
	_, err = stmt.Exec(a.Id, a.Name, a.Detail, a.Coverimage)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Fatal("insert success")
	}
	c.IndentedJSON(http.StatusCreated, a)
}
