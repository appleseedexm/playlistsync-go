package tidal2

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"playlistsync/main/spotify"
	"playlistsync/main/util"
	"strconv"
	"strings"
)


type TidalWebApi struct {
	username    string
	password    string
	token       string
	quality     string
	countryCode string
	bearerToken string
	sessionId   string
	origin      string
}

type Song struct {
	artist   string
	songName string
}

func SyncSongs(songs []spotify.Song, envVars util.EnvVars) {

	tidal := TidalWebApi{
		username:    "",
		password:    "",
		countryCode: "CH",
		bearerToken: "Bearer " + envVars.TidalBearerToken,
	}

	var songIds []string

	for _, song := range songs {
		songId := tidal.getSongId(song.Artist, song.SongName)
		if songId == 0 {
			continue
		}
		songIds = append(songIds, strconv.FormatFloat(songId, 'f', -1, 64))
	}
	fmt.Println("Adding songs to playlist:")
	fmt.Println(songIds)

	newPlaylistUuid := tidal.createPlaylist("APITESTttfasdf")
	etag := tidal.getPlaylist(newPlaylistUuid)
	tidal.addSongToPlaylist(newPlaylistUuid, songIds, etag)

}

func (h *TidalWebApi) getSongId(artist string, songName string) float64 {
	type SearchResult struct {
		Tracks struct {
			Items []struct {
				Id float64
			}
		}
	}

	fmt.Println("#########################################################")
	fmt.Println("Searching for song:")
	fmt.Println(artist)
	fmt.Println(songName)

	artist = strings.ReplaceAll(artist, " ", "+")
	songName = strings.ReplaceAll(songName, " ", "+")
	songName = strings.ReplaceAll(songName, "’", "")
	url := fmt.Sprintf("https://listen.tidal.com/v2/search/?includeContributors=true&includeDidYouMean=true&includeUserPlaylists=false&limit=25&query=%s+%s&supportsUserData=true&types=TRACKS&countryCode=CH&locale=en_US&deviceType=BROWSER", artist, songName)
	//url = "https://listen.tidal.com/search?q=the%20kids%20arent%20alright"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", h.bearerToken)
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)

	fmt.Println("Response Status:", res.Status)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var searchResult SearchResult
	err = json.Unmarshal(body, &searchResult)
	if err != nil {
		panic(err)
	}

	if len(searchResult.Tracks.Items) == 0 {

		fmt.Println("#########################################################")
		fmt.Println("#########################################################")
		fmt.Println("#########################################################")
		fmt.Println("#########################################################")
		fmt.Println("#########################################################")
		fmt.Println("#########################################################")
		fmt.Println("#########################################################")
		fmt.Println("No song found")
		fmt.Println(artist)
		fmt.Println(songName)
		return 0
	}

	return searchResult.Tracks.Items[0].Id
}

func (h *TidalWebApi) getPlaylist(playlistId string) string {
	//https://tidal.com/browse/playlist/14ab0fc9-26c8-44fb-9aa3-b07b539a4b0e
	//TidalAPI.prototype.getPlaylists = function (user, callback) {
	//var self = this;
	//self._baseRequest('/users/' + (user.id || user) + "/playlists", {
	//limit: user.limit || 999,
	//offset: user.offset || 0,
	//countryCode: _countryCode
	//}, 'userPlaylists', callback);
	//}

	//https://listen.tidal.com/v1/playlists/c8c9d861-971c-4bc9-ab7d-85d7450c732c?countryCode=CH&locale=en_US&deviceType=BROWSER

	// POST song to playlist
	// onArtifactNotFound=FAIL&onDupes=FAIL&trackIds=236360293
	// https://listen.tidal.com/v1/playlists/46f7f50b-829c-4aa7-b936-44f67209befe/items?countryCode=CH&locale=en_US&deviceType=BROWSER
	playlistUrl := "https://listen.tidal.com/v1/playlists/" + playlistId + "?countryCode=CH&locale=en_US&deviceType=BROWSER"

	req, err := http.NewRequest("GET", playlistUrl, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", h.bearerToken)
	res, err := http.DefaultClient.Do(req)

	fmt.Println("Response Status:", res.Status)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	if res.StatusCode != 200 {
		panic("Error")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))

	return res.Header.Get("ETag")

}

func (h *TidalWebApi) createPlaylist(name string) string {

	type TidalPlaylist struct {
		Data struct {
			Uuid string
		}
	}

	url := fmt.Sprintf("https://listen.tidal.com/v2/my-collection/playlists/folders/create-playlist?description=&folderId=root&isPublic=true&name=%s&countryCode=CH&locale=en_US&deviceType=BROWSER", name)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", h.bearerToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var playlist TidalPlaylist
	err2 := json.Unmarshal(body, &playlist)
	if err2 != nil {
		panic(err2)
	}

	playlistId := playlist.Data.Uuid

	return playlistId
}

func (h *TidalWebApi) addSongToPlaylist(playlistId string, songIds []string, etag string) {

	// POST song to playlist
	// onArtifactNotFound=FAIL&onDupes=FAIL&trackIds=236360293
	// https://listen.tidal.com/v1/playlists/46f7f50b-829c-4aa7-b936-44f67209befe/items?countryCode=CH&locale=en_US&deviceType=BROWSER

	values := url.Values{
		"onArtifactNotFound": {"FAIL"},
		"onDupes":            {"FAIL"},
		"trackIds":           {strings.Join(songIds, ",")},
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("https://listen.tidal.com/v1/playlists/%s/items?countryCode=CH&locale=en_US&deviceType=BROWSER", playlistId), strings.NewReader(values.Encode()))

	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", h.bearerToken)
	req.Header.Set("If-None-Match", etag)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	//if res.StatusCode != 200 {
	//panic("Error")
	//}

	_, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

}
