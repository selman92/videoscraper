package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("Starting scraping")
	startUrl := "https://mywebsite.net/"

	c := colly.NewCollector(

		colly.AllowedDomains("mysite.net", ""),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if !strings.Contains(link, "watch") {
			return
		}

		// Print link
		fmt.Printf("Link found: %s\n", link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML("source", func(e *colly.HTMLElement) {
		videoLink := e.Attr("src")

		go DownloadFile(videoLink, "videos")
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(response *colly.Response) {

	})
	err := c.Visit(startUrl)

	if err != nil {
		fmt.Println(err)
	}

}

func DownloadFile(url, folder string) error {
	// Generate a GUID for the filename
	fileName := uuid.New().String() + ".mp4"
	filePath := filepath.Join(folder, fileName)

	// Create the file where the MP4 will be saved
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	// Get the data from the URL
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	// Check if the HTTP status code is OK (200)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the data to the file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	fmt.Printf("File saved to: %s\n", filePath)
	return nil
}
