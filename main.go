package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/mmcdole/gofeed"
)

const settingsPath = "/.config/rss-to-pocket/settings.yaml"
const credentialsPath = "/.config/rss-to-pocket/credentials.yaml"

type RssFeed struct {
	Url                 string `yaml:"url"`
	MaxNumberOfArticles int    `yaml:"max_number_of_articles"`
}

type Settings struct {
	RssFeeds []RssFeed `yaml:"rss_feeds"`
}
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
		for _, item := range feed.Items {
			fmt.Println("")
			time.Sleep(2 * time.Second)
			fmt.Println(item.Link)
			fmt.Println(item.PublishedParsed)
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
