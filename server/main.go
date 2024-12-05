package main

import (
	"database/sql"
	"log"
	"net/http"
	"sync"
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

func getDatabaseHandle() *sql.DB {
	db, err := sql.Open("sqlite", "file:./top2000.sqlite3?mode=ro")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func initApplication() *sql.DB {
	db := getDatabaseHandle()
	return db
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
	db := initApplication()
	defer db.Close()
	randomTrack = getRandomTrack(db)
	indexTemplate, err := template.ParseFiles("templates/index.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
