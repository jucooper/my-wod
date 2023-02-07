package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli"
)

// PushJerkEndpoint is the url for PushJerk's workout
const PushJerkEndpoint = "http://pushjerk.com/index.php/workout/%s-%s-%d-%d/"

// PushJerkScrape scrapes pushjerk.com
func PushJerkScrape(c *cli.Context, date time.Time) {
	fmt.Println("Push Jerk WOD")
	url := fmt.Sprintf(PushJerkEndpoint,
		Abbreviate(date.Weekday().String()),
		Abbreviate(date.Month().String()),
		date.Day(),
		date.Year(),
	)
	fmt.Println(url)
	content := Scrape(url, "", ".entry-content")
	fmt.Println(content)
}

// Loader handles the request and loads the html
func Loader(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode == 404 {
		fmt.Println("Workout is not available.")
		os.Exit(404)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc, err
}

// Scrape scrapes the html page for the specified element
func Scrape(url string, date string, element string) string {
	doc, _ := Loader(url)
	s := doc.Find(element).First()
	text := s.Find("p").Text()
	return text
}

// Abbreviate modifies the weekday/month date input to follow PushJerk's url format
func Abbreviate(date string) string {
	return strings.ToLower(date[0:3])
}

func main() {
	app := cli.NewApp()
	app.Name = "my-wod"
	app.Version = "1.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name: "Justin Cooper",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "today",
			Usage: "Display today's workout.",
			Action: func(c *cli.Context) error {
				PushJerkScrape(c, time.Now())
				return nil
			},
		},
		{
			Name:  "tomorrow",
			Usage: "Display tomorrow's workout.",
			Action: func(c *cli.Context) error {
				PushJerkScrape(c, time.Now().AddDate(0, 0, 1))
				return nil
			},
		},
		{
			Name:  "yesterday",
			Usage: "Display yesterday's workout.",
			Action: func(c *cli.Context) error {
				PushJerkScrape(c, time.Now().AddDate(0, 0, -1))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
