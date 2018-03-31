package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli"
)

// InvictusEndpoint is the url for Invictus' performance workout
const InvictusEndpoint = "http://www.crossfitinvictus.com/wod/"

// PushJerkEndpoint is the url for PushJerk's workout
const PushJerkEndpoint = "http://pushjerk.com/index.php/workout/"

// InvictusScrape scrapes crossfitinvictus.com
func InvictusScrape() {
	fmt.Println("Invictus Fitness WOD")
	content := Scrape(InvictusEndpoint, ".entry-content")
	fmt.Println(content)
}

// PushJerkScrape scrapes rushjerk.com
func PushJerkScrape() {
	fmt.Println("Push Jerk WOD")
	content := Scrape(PushJerkEndpoint, ".entry-content")
	fmt.Println(content)
}

// Loader handles the request and loads the html
func Loader(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc, err
}

// Scrape scrapes the url and element
func Scrape(url string, element string) string {
	doc, _ := Loader(url)
	s := doc.Find(".entry-content").First()
	text := s.Find("p").Text()
	return text
}

func main() {
	app := cli.NewApp()
	app.Name = "my-wod"
	app.Version = "1.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Justin Cooper",
			Email: "justinwcooper@outlook.com",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "today",
			Aliases: []string{"t"},
			Usage:   "Display today's workout.",
			Action: func(c *cli.Context) error {
				InvictusScrape()
				fmt.Print("\n\n")
				PushJerkScrape()

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
