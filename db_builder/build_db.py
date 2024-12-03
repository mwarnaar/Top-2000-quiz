import json
import sqlite3

con = sqlite3.connect("top2000.sqlite3")
cursor = con.cursor()

# CREATE TABLE all_songs(track_id VARCHAR(255), track_name VARCHAR(255), track_href VARCHAR(1024), artist_id VARCHAR(255), artist_name VARCHAR(255), album_id VARCHAR(255), album_name VARCHAR(255), album_release_date VARCHAR(10), album_image_url VARCHAR(1024));

cursor.execute("DELETE FROM all_songs")

with open("tracks_dump.json", "rb") as tracks:
    tracks = json.loads(tracks.read())
    for track in tracks:
        track = track.get("track")

        track_id = track.get("id")
        track_name = track.get("name")
        track_href = track.get("uri")

        artist_name = track.get("artists")[0].get("name")
        artist_id = track.get("artists")[0].get("id")

        album_release_date = track.get("album").get("release_date").replace("-", "")
        album_id = track.get("album").get("id")
        album_name = track.get("album").get("name")
        album_image_url = track.get("album").get("images")[0].get("url")
        print(
            artist_name,
            artist_id,
            track_name,
            track_id,
            album_release_date,
            album_image_url,
        )
        record = (
            track_id,
            track_name,
            track_href,
            artist_id,
            artist_name,
            album_id,
            album_name,
            album_release_date,
            album_image_url,
        )
        sql = "INSERT INTO all_songs(track_id, track_name, track_href, artist_id, artist_name, album_id, album_name, album_release_date, album_image_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
        cursor.execute(sql, record)
con.commit()
