package main

import (
	"database/sql"
	"log"
	"net/http"
	"sync/atomic"
	"text/template"

	_ "modernc.org/sqlite"
)

type track struct {
	Track_id           string
	Track_href         string
	Artist_name        string
	Track_name         string
	Album_name         string
	Album_image_url    string
	Album_release_date string
}

// Truncate release date to year (Y)
func getAlbumReleaseYear(rawInput string) string {
	return rawInput[0:4]
}

// Get random track from SQLite
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
		&t.Track_id,
		&t.Track_href,
		&t.Artist_name,
		&t.Track_name,
		&t.Album_name,
		&t.Album_image_url,
		&t.Album_release_date,
	)
	if err != nil {
		log.Fatal(err)
	}
	t.Album_release_date = getAlbumReleaseYear(t.Album_release_date)

	return t
}

// Open SQLite database
func getDatabaseHandle() *sql.DB {
	db, err := sql.Open("sqlite", "file:./top2000.sqlite3?mode=ro")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// State about running DB query job and current provided random track
var (
	jobIsRunning int32
	randomTrack  track
)

func getNewRandomTrackIfNotLoading(db *sql.DB) {
	if atomic.CompareAndSwapInt32(&jobIsRunning, 0, 1) {
		go func() {
			randomTrack = getRandomTrack(db)
			atomic.StoreInt32(&jobIsRunning, 0)
		}()
	}
}

func main() {
	db := getDatabaseHandle()
	defer db.Close()
	randomTrack = getRandomTrack(db)
	indexTemplate, err := template.ParseFiles("templates/index.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		err = indexTemplate.Execute(w, randomTrack)
		if err != nil {

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Fatal(w)
			return
		}

		getNewRandomTrackIfNotLoading(db)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
