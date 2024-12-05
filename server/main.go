package main

import (
	"database/sql"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
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
		WHERE rowid = (abs(random()) % 1996) + 1
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
	db, err := sql.Open("sqlite", "file:./top2000.sqlite3?mode=ro")
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

var (
	jobIsRunning      bool
	jobIsRunningMutex sync.Mutex
	randomTrack       track
)

func getNewRandomTrackIfNotLoading(db *sql.DB) {
	jobIsRunningMutex.Lock()
	start := !jobIsRunning
	jobIsRunning = true
	jobIsRunningMutex.Unlock()
	if start {
		go func() {
			jobIsRunningMutex.Lock()
			randomTrack = getRandomTrack(db)
			jobIsRunning = false
			jobIsRunningMutex.Unlock()
		}()
	}
}

func main() {
	router, db := initApplication()
	defer db.Close()
	randomTrack = getRandomTrack(db)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"track_id":           randomTrack.track_id,
			"track_uri":          randomTrack.track_href,
			"artist_name":        randomTrack.artist_name,
			"album_name":         randomTrack.album_name,
			"track_name":         randomTrack.track_name,
			"album_image_url":    randomTrack.album_image_url,
			"album_release_date": randomTrack.album_release_date,
		})
		getNewRandomTrackIfNotLoading(db)
	})
	router.Run(":8080")
}
