package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/mmcdole/gofeed"
)

const settingsPath = "/.config/rss-to-pocket/settings.yaml"
const credentialsPath = "/.config/rss-to-pocket/credentials.yaml"
const savedEntriesPath = "/.config/rss-to-pocket/savedEntries.csv"
const singleAddUrl = "https://getpocket.com/v3/add"

/** Data structures in {@link settingsPath} */
type RssFeed struct {
	Url                 string `yaml:"url"`
	MaxNumberOfArticles int    `yaml:"max_number_of_articles"`
	Tag                 string `yaml:"tag"`
}
type Settings struct {
	RssFeeds []RssFeed `yaml:"rss_feeds"`
}

/** Data structure in {@link credentialsPath} */
type Credentials struct {
	ConsumerKey string `yaml:"consumer_key"`
	AccessToken string `yaml:"access_token"`
}

func main() {
	settings := getSettings()
	credentials := getCredentials()

	if credentials.AccessToken == "" {
		Authenticate(credentials)

	}
	for _, settingEntry := range settings.RssFeeds {
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(settingEntry.Url)
		// Assume these are sorted time desc
		count := 0
		max := 5
		if settingEntry.MaxNumberOfArticles != 0 {
			max = settingEntry.MaxNumberOfArticles
		}
		for _, item := range feed.Items {
			if count > max {
				fmt.Print("Max articles reached")
				break
			}
			fmt.Println("")
			time.Sleep(30 * time.Second)
			addItemToPocket(item, credentials, settingEntry.Tag)
			fmt.Println(item.Link)
			fmt.Println(item.PublishedParsed)
			count++
		}

	}
}
func getSettings() Settings {
	settings := Settings{}

	path := getPathFromHome(settingsPath)
	content := getFileContents(path)
	err := yaml.Unmarshal(content, &settings)
	if err != nil {
		log.Fatalf("unmarshall err %s", err)
	}
	return settings

}
func getCredentials() Credentials {
	credentials := Credentials{}

	path := getPathFromHome("/.config/rss-to-pocket/credentials.yaml")
	content := getFileContents(path)
	err := yaml.Unmarshal(content, &credentials)
	if err != nil {
		log.Fatalf("unmarshall err %s", err)
	}
	return credentials

}

func getFileContents(path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return content
}

func getPathFromHome(path string) string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Can't find home directory")
	}
	return homePath + path
}
func addItemToPocket(item *gofeed.Item, credentials Credentials, tag string) {
	postBody := fmt.Sprintf(`{"consumer_key": "%s", "access_token": "%s", "url":"%s", "tags": "%s"}`,
		credentials.ConsumerKey,
		credentials.AccessToken,
		item.Link,
		tag)

	_, responseCode := makePostRequest(singleAddUrl, postBody)
	writeTime := time.Now().String()
	if strings.Contains(responseCode, "200") {
		f, err := os.OpenFile(getPathFromHome(savedEntriesPath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write([]byte(writeTime + "," + item.Link + "\n")); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

func makePostRequest(url string, postBody string) (string, string) {
	jsonBody := []byte(postBody)
	bodyReader := bytes.NewReader(jsonBody)

	req, _ := http.NewRequest(http.MethodPost, url, bodyReader)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	// fmt.Println("")
	// fmt.Printf("Request Body: %s\n", postBody)
	// fmt.Printf("Request URL %s\n", url)
	// fmt.Printf("Response Body: %s\n", string(body))
	// fmt.Printf("Response Headers: %s\n", response.Header)
	// fmt.Println("")
	return string(body), response.Status
}
