package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"playlistsync/main/auth"

	"github.com/rapito/go-spotify/spotify"
	"golang.org/x/oauth2"
)

type OAuthRedirectHandler struct {
	State        string
	CodeVerifier string
	OAuthConfig  *oauth2.Config
}

func main() {
	fmt.Println("Hello, World!")

	clientId := ""

	state := auth.Authorize(clientId, "http://localhost:8080")

	oauth := &OAuthRedirectHandler{
		State:        state,
		CodeVerifier: "code_verifier",
		OAuthConfig: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: "",
			RedirectURL:  "http://localhost:8080",
			Scopes:       []string{"r_usr", "w_usr"},
		},
	}


    http.HandleFunc("/", oauth.startListeningServer)
    log.Fatal(http.ListenAndServe(":8080", nil))

}

func (h *OAuthRedirectHandler) startListeningServer(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	state := query.Get("state")


    // todo: implement state
	if subtle.ConstantTimeCompare([]byte(h.State), []byte(state)) == 0  || true{
		http.Error(rw, "Invalid state", http.StatusBadRequest)
		return
	}

	code := query.Get("code")

	if code == "" {
		http.Error(rw, "No code", http.StatusBadRequest)
		return
	}

	token, err := h.OAuthConfig.Exchange(req.Context(), code, oauth2.SetAuthURLParam("code_verifier", h.CodeVerifier))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(token)

	//server := &http.Server{Addr: redirectUrl}
	//_ = server.ListenAndServe()

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
