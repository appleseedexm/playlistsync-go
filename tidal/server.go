package tidal

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"playlistsync/main/util"
)

type OAuthRedirectHandler struct {
	State        string
	CodeVerifier string
	ClientId     string
	Code         string
	authUrl      string
	redirectUri  string
}

func Server(){

	clientId := ""

	state, code_verifier := Authorize(clientId, "http://localhost:8080/authorized")
	oauth := &OAuthRedirectHandler{
		State:        state,
		CodeVerifier: code_verifier,
		ClientId:     clientId,
		Code:         "",
		authUrl:      "https://auth.tidal.com/v1/oauth2/token",
		redirectUri:  "http://localhost:8080/token",
	}

	http.HandleFunc("/authorized", oauth.listenAuthorization)
	http.HandleFunc("/token", oauth.listenToken)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (h *OAuthRedirectHandler) listenToken(rw http.ResponseWriter, req *http.Request) {

    query := req.URL.Query()

    accessToken := query.Get("access_token")
    tokenType := query.Get("token_type")
    scope := query.Get("scope")
    expiresIn := query.Get("expires_in")
    refreshToken := query.Get("refresh_token")
    
    fmt.Println(accessToken)
    fmt.Println(tokenType)
    fmt.Println(scope)
    fmt.Println(expiresIn)
    fmt.Println(refreshToken)

}

func (h *OAuthRedirectHandler) listenAuthorization(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	state := query.Get("state")

	// todo: implement state
	if subtle.ConstantTimeCompare([]byte(h.State), []byte(state)) == 1 {
		http.Error(rw, "Invalid state", http.StatusBadRequest)
	}

	code := query.Get("code")

	if code == "" {
		http.Error(rw, "No code", http.StatusBadRequest)
		return
	}

	h.Code = code

	//curl -X POST \
	//-d "grant_type=authorization_code" \
	//-d "client_id=<CLIENT_ID>" \
	//-d "code=<CODE>" \
	//-d "redirect_uri=<REDIRECT_URI>" \
	//-d "code_verifier=<CODE_VERIFIER>" \
	//"https://auth.tidal.com/v1/oauth2/token"

    res, err := http.PostForm("https://auth.tidal.com/v1/oauth2/token", url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {h.ClientId},
		"code":          {h.Code},
		"redirect_uri":  {h.redirectUri},
		"code_verifier": {h.CodeVerifier},
	})

    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }

	fmt.Println(res.Body.Close().Error())

	//token, err := h.OAuthConfig.Exchange(req.Context(), code, oauth2.SetAuthURLParam("code_verifier", h.CodeVerifier))

	//if err != nil {
	//http.Error(rw, err.Error(), http.StatusInternalServerError)
	//return
	//}

	//fmt.Println(token)

	//server := &http.Server{Addr: redirectUrl}
	//_ = server.ListenAndServe()

}

func Authorize(clientId string, authUrl string) (string, string){

	//server := http.Server{}

	//scope := "collection.read,collection.write,playlists.read,playlists.write"
    scope := "collection.read"

	code_challenge := util.RandomBytesInBase64(int(32))
	println(code_challenge)
	s2 := sha256.New()
	s2.Write([]byte(code_challenge))
	code_challenge_hashed_base64 := base64.RawURLEncoding.EncodeToString(s2.Sum(nil))
	println(code_challenge_hashed_base64)

    state := util.RandomBytesInHex(int(24))

	authorizationUrl := fmt.Sprintf("https://login.tidal.com/authorize"+
		"?response_type=code"+
		"&client_id=%s"+
		"&redirect_uri=%s"+
		"&scope=%s"+
		"&code_challenge_method=S256"+
		"&code_challenge=%s"+
        "&state=%s", clientId, authUrl, scope, code_challenge_hashed_base64, state)

	fmt.Println(authorizationUrl)

    return state, code_challenge

	//http.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {

	//})
}
