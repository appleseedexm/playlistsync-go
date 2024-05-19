package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

func Authorize(clientId string, authUrl string) {

	//server := http.Server{}

	scope := "r_usr,w_usr"

	code_challenge := randomBytesInBase64(int(32))
	s2 := sha256.New()
	s2.Write([]byte(code_challenge))
	code_challenge_hashed_base64 := base64.RawURLEncoding.EncodeToString(s2.Sum(nil))

	authorizationUrl := fmt.Sprintf("https://login.tidal.com/authorize"+
		"?response_type=code"+
		"&client_id=%s"+
		"&redirect_uri=%s"+
		"&scope=%s"+
		"&code_challenge_method=S256"+
		"&code_challenge=%s", clientId, authUrl, scope, code_challenge_hashed_base64)

	fmt.Println(authorizationUrl)

	//http.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {

	//})
}

func randomBytesInBase64(count int) string {
	buf := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(buf)
}
