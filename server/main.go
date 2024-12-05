package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type track struct {
	track_id           string
	track_href         string
	artist_name        string
	track_name         string
	album_name         string
	album_image_url    string
	album_release_date string
}

// Truncate release date YMD, YM, Y strings to Y
func getAlbumReleaseYear(rawInput string) string {
	return rawInput[0:4]
}

func getRandomTrack(db *sql.DB) track {
	row := db.QueryRow(`
		SELECT
			track_id,
			track_href,
			artist_name,
			track_name,
			album_name,
			album_image_url,
			album_release_date
		FROM all_songs
		ORDER BY RANDOM()
		LIMIT 1;
	`)
	var t track
	err := row.Scan(
		&t.track_id,
		&t.track_href,
		&t.artist_name,
		&t.track_name,
		&t.album_name,
		&t.album_image_url,
		&t.album_release_date,
	)
	if err != nil {
		log.Fatal(err)
	}
	t.album_release_date = getAlbumReleaseYear(t.album_release_date)

	return t
}

func getDatabaseHandle() *sql.DB {
	db, err := sql.Open("sqlite3", "./top2000.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func initApplication() (*gin.Engine, *sql.DB) {
	db := getDatabaseHandle()
	router := gin.Default()
	router.LoadHTMLFiles("templates/index.tmpl")
	router.Static("/public", "./public")
	return router, db
}

func main() {
	router, db := initApplication()
	defer db.Close()

	router.GET("/", func(c *gin.Context) {
		track := getRandomTrack(db)
		fmt.Printf("URI: %s, artist: %s, title: %s\n", track.track_href, track.artist_name, track.track_name)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"track_id":           track.track_id,
			"track_uri":          track.track_href,
			"artist_name":        track.artist_name,
			"album_name":         track.album_name,
			"track_name":         track.track_name,
			"album_image_url":    track.album_image_url,
			"album_release_date": track.album_release_date,
		})
	})
	router.Run(":8080")
}
