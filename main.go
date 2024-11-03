package main

import (
	"fmt"
	"os"
	"playlistsync/main/spotify"
	"playlistsync/main/tidal2"
	"playlistsync/main/util"

	"github.com/joho/godotenv"
)

func main() {
    env := load_dotenv()

	spotifySongs := spotify.GetPlaylistFromSpotify(env)
	fmt.Println(`Adding songs to playlist:`)
	fmt.Println(spotifySongs)

	tidal2.SyncSongs(spotifySongs, env)

}

func load_dotenv() util.EnvVars {
    fmt.Sprint("Loading environment variables..")

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Sprint("No .env found...")
	}

	return util.EnvVars{
		SpotifyClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		TidalBearerToken:    os.Getenv("TIDAL_BEARER_TOKEN")}

}
