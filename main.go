package main

import (
	
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	"crypto/sha1"
	"encoding/hex"

	"github.com/gin-contrib/cors"
	
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	

)

type FormData struct {
	IpAddress  string `json:"ipaddress"`
	Image      string `json:"image"`
}
type IP struct {
	IpAddress  string `json:"ipaddress"`
}

func main() {
	var router = gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))
	router.LoadHTMLGlob("templates/*.html")
	router.GET("/get_form", GetFormHandler)
	router.GET("/", GetFormHandler)
	router.POST("/save_form", FormHandler)
	router.POST("/get_image_ip", ImageHandler) 
	router.GET("/get_image", GetImageHandler) 
	router.Run(":8080")
}
func GetImageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "image.html", gin.H{})
}
func GetFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{})
}
func FormHandler(c *gin.Context) {
	var data FormData
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err)
	}
	date := time.Now()	
	
	h := sha1.New()
	h.Write([]byte(data.Image))
	hash1 := h.Sum(nil)
	
	// md5HashInBytes := md5.Sum([]byte("Sum returns bytes"))
	// md5HashInString := hex.EncodeToString(md5HashInBytes[:])
	// fmt.Println(md5HashInString)
	fmt.Println("HHHHHHHHHHH")
	fmt.Println(hex.EncodeToString(hash1))

	db := DBconnect()
	defer db.Close()
	_, err = db.Exec("insert into public.imagesHashed(ipaddress_user,image,created_date,hashimage) values($1,$2,$3,$4)", data.IpAddress, data.Image, date, hex.EncodeToString(hash1))
	
	if err != nil {
		fmt.Println("here")
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "data saved",
	})
}

 func ImageHandler(c *gin.Context) {
	var data IP
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err)
	}

	db := DBconnect()
	defer db.Close()
	var image string
	fmt.Println(data.IpAddress)
	res, err := db.Query("select image from public.imagesHashed where ipaddress_user=$1 and created_date=(select max(created_date) from public.imagesHashed where ipaddress_user=$1)", data.IpAddress)
	if err != nil {
		fmt.Println(err)
	}
	for res.Next() {
		err = res.Scan(&image)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("success")
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"image":  image,
	})
} 
func DBconnect() *sql.DB {

	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=roadmakerDB sslmode=disable")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(db)
	// close database

	// check db
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")
	return db
}
