import spotipy
from spotipy.oauth2 import SpotifyClientCredentials
import json
from dotenv import load_dotenv
from os import environ

load_dotenv()  # take environment variables from .env.

sp = spotipy.Spotify(
    auth_manager=SpotifyClientCredentials(
        client_id=environ.get("SPOTIFY_CLIENT_ID"),
        client_secret=environ.get("SPOTIFY_SECRET_ID"),
    )
)

results = sp.playlist_items(playlist_id="1DTzz7Nh2rJBnyFbjsH1Mh")

tracks = results.get("items")
while results["next"]:
    results = sp.next(results)
    tracks.extend(results["items"])

print(json.dumps(tracks))
