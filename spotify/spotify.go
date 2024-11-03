package spotify

import (
	"encoding/json"
	"github.com/rapito/go-spotify/spotify"
	"playlistsync/main/util"
)

type Song struct {
	Artist   string
	SongName string
}

func GetPlaylistFromSpotify(envVars util.EnvVars) []Song {

	clientId := envVars.SpotifyClientId
	clientSecret := envVars.SpotifyClientSecret
	spot := spotify.New(clientId, clientSecret)

	result, _ := spot.Get("playlists/%s", nil, "5EGiNnE8oWvzVpHnAZVF3O")

	var playlistResponse PlaylistResponse

	json.Unmarshal([]byte(result), &playlistResponse)

	var songs []Song

	for _, track := range playlistResponse.AllTracks.TrackWithMetaData {
		songs = append(songs, Song{
			Artist:   track.Track.Artists[0].Name,
			SongName: track.Track.Name,
		})
	}

	return songs
}

type PlaylistResponse struct {
	AllTracks AllTracks `json:"tracks"`
}

type AllTracks struct {
	TrackWithMetaData []TrackMetaData `json:"items"`
}

type TrackMetaData struct {
	Track Track `json:"track"`
}

type Track struct {
	Name    string   `json:"name"`
	Artists []Artist `json:"artists"`
}

type Artist struct {
	Name string `json:"name"`
}
