package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const redirectUri = "http://localhost:3333/rss-to-pocket/auth"

// check yaml

// pubic method, run main auth flow
// https://getpocket.com/developer/docs/authentication
func Authenticate(credentials Credentials) {
	fmt.Println("Starting Server")
	requestToken := getRequestToken(credentials.ConsumerKey)
	go runServer(credentials, requestToken)

	fmt.Println("")
	fmt.Printf("Authenticate with pocket at:")
	fmt.Printf("https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s", requestToken, redirectUri)
}
// Run a local http server to handle Pocket callback
func runServer(credentials Credentials, requestToken string) {
	http.HandleFunc("/rss-to-pocket/auth", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Getting Pocket Access token")
		accessToken := getPocketAccessToken(credentials.ConsumerKey, requestToken)
		writeAccessTokenToYaml(accessToken, credentials)
	})

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getRequestToken(consumerKey string) string {
	fmt.Println()
	postBody := fmt.Sprintf(`{"consumer_key": "%s", "redirect_uri": "%s"}`, consumerKey, "localhost:3333")
	body, _ := makePostRequest("https://getpocket.com/v3/oauth/request", postBody)

	_, code, found := strings.Cut(body, "code=")
	if !found {
		log.Fatal("Can't find code")
	}
	return code

}
func getPocketAccessToken(consumerKey string, requestToken string) string {
	postBody := fmt.Sprintf(`{"consumer_key": "%s", "code": "%s"}`, consumerKey, requestToken)
	body, _ := makePostRequest("https://getpocket.com/v3/oauth/authorize", postBody)
	fmt.Println(body)
	_, after, found := strings.Cut(body, "access_token=")
	if !found {
		log.Fatal("Can't find code")
	}
	code, _, found := strings.Cut(after, "&username=")
	if !found {
		log.Fatal("Can't find code")
	}

	return code
}
func writeAccessTokenToYaml(accessToken string, credentials Credentials) {
	fmt.Println("")
	fmt.Printf("Writing access token %s to file", accessToken)

	credentials.AccessToken = accessToken
	data, err := yaml.Marshal(&credentials)

	if err != nil {
		log.Fatal(err)
	}

	err2 := os.WriteFile(getPathFromHome(credentialsPath), data, 0)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("Access token written")
}
