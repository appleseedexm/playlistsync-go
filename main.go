package main

import (
	"encoding/json"
	"fmt"
	"github.com/rapito/go-spotify/spotify"
)

type OAuthRedirectHandler struct {
	State        string
	CodeVerifier string
	ClientId     string
	Code         string
	authUrl      string
	redirectUri  string
}

func main() {
	fmt.Println("Hello, World!")


}

func getPlaylistFromSpotify() {

	spot := spotify.New("", "")
	spot.Authorize()

	result, _ := spot.Get("playlists/%s", nil, "5QIZGd9DfVRrcDJJPBISQv")

	var playlistResponse PlaylistResponse

	error := json.Unmarshal([]byte(result), &playlistResponse)
	fmt.Println(error)

	fmt.Println(string((playlistResponse.AllTracks.TrackWithMetaData[0].Track.Artists[0].Name)))
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
