Download Spotify playlist tracks to SQLite
==========================================
Build an SQLite database with Spotify songs from a specified playlist.

* Copy `.env.example` to `.env`, specify your Spotify Client ID and Secret.
* Create an sqlite DB file and table, see `create_db.sql`.
* Fetch and fill the database:
```shell
python3 -m venv venv
source venv/bin/activate

python3 -m pip install -r requirements.txt
python3 fetch_playlist.py > tracks_dump.json
python3 build_db.py
```

You should now have a `top2000.sqlite3` database file that can be used as tracks source for the server (`../server/top2000.sqlite3`).

