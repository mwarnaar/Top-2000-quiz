package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./top2000.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	router.LoadHTMLFiles("templates/index.tmpl")
	router.Static("/public", "./public")

	router.GET("/", func(c *gin.Context) {
		row := db.QueryRow("SELECT track_id, track_href, artist_name, track_name, album_name, album_image_url, album_release_date FROM all_songs ORDER BY RANDOM() LIMIT 1;")
		var track_id string
		var track_href string
		var artist_name string
		var track_name string
		var album_name string
		var album_image_url string
		var album_release_date string
		err = row.Scan(&track_id, &track_href, &artist_name, &track_name, &album_name, &album_image_url, &album_release_date)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(track_href, artist_name, track_name)
		var album_release_date_result string
		if len(album_release_date) >= 8 {
			year := album_release_date[0:4]
			month := album_release_date[4:6]
			day := album_release_date[6:8]
			album_release_date_result = fmt.Sprintf("%s-%s-%s", year, month, day)
		} else if len(album_release_date) >= 6 {
			year := album_release_date[0:4]
			month := album_release_date[4:6]
			album_release_date_result = fmt.Sprintf("%s-%s", year, month)
		} else {
			album_release_date_result = fmt.Sprintf("%s", album_release_date)
		}
		if err != nil {
			log.Fatal(err)
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"track_id":           track_id,
			"track_uri":          track_href,
			"artist_name":        artist_name,
			"album_name":         album_name,
			"track_name":         track_name,
			"album_image_url":    album_image_url,
			"album_release_date": album_release_date_result,
		})
	})
	router.Run(":8080")
}
