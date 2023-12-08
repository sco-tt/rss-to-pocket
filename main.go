package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/mmcdole/gofeed"
)

type RssFeed struct {
	Url                 string `yaml:"url"`
	MaxNumberOfArticles int    `yaml:"max_number_of_articles"`
}

type Settings struct {
	PocketApiConsumerKey string    `yaml:"pocket_api_consumer_key"`
	RssFeeds             []RssFeed `yaml:"rss_feeds"`
}

func main() {

	settings := getSettings()
	for _, settingEntry := range settings.RssFeeds {
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(settingEntry.Url)
		// Assume these are sorted time desc
		for _, item := range feed.Items {
			time.Sleep(2 * time.Second)
			fmt.Println(item.Link)
			fmt.Println(item.PublishedParsed)
		}

	}
}
func getSettings() Settings {
	settings := Settings{}
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Can't find home directory")
	}

	path := homePath + "/.config/rss-to-pocket/settings.yaml"
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(content, &settings)
	if err != nil {
		log.Fatalf("unmarshall err %s", err)
	}
	return settings

}
