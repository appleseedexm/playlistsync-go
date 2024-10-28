package main

import (
	"fmt"
	"playlistsync/main/tidal2"
    "playlistsync/main/spotify"
)

func main() {
	fmt.Println("Hello, World!")

    spotifySongs := spotify.GetPlaylistFromSpotify()
    fmt.Println(`Adding songs to playlist:`)
    fmt.Println(spotifySongs)

	tidal2.Serve(spotifySongs)

}
