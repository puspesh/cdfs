package models

import (
	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/drive/v2"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	FOLDER_NAME = "cdfs"
	DESCRIPTION = "cdfs data"
)

// Settings for authorization.
var config = &oauth.Config{
	ClientId:     "242859917159-9q0mcad9cuhflb42g7dq5lgl46tccsm6.apps.googleusercontent.com",
	ClientSecret: "fCmLibbJ-p-28AW3HN1NiC_e",
	Scope:        "https://www.googleapis.com/auth/drive",
	RedirectURL:  "http://localhost/oauth2callback",
	AuthURL:      "https://accounts.google.com/o/oauth2/auth",
	TokenURL:     "https://accounts.google.com/o/oauth2/token",
}

// Uploads a file to Google Drive
func (gd *GoogleDrive) Upload(file string) string {
	// Generate a URL to visit for authorization.
    var c oauth.CacheFile = oauth.CacheFile(gd.authString)
    tok,_ := c.Token()
	t := &oauth.Transport{
		Config:    config,
		Transport: http.DefaultTransport,
        Token: tok,
	}
	// Create a new authorized Drive client.
	svc, err := drive.New(t.Client())
	if err != nil {
		fmt.Printf("An error occurred creating Drive client: %v\n", err)
	}

	// Define the metadata for the file we are going to create.
	f := &drive.File{
		Title:       FOLDER_NAME,
		Description: DESCRIPTION,
	}

	// Read the file data that we are going to upload.
	m, err := os.Open(file)
	if err != nil {
		fmt.Printf("An error occurred reading the document: %v\n", err)
	}

	tStart := time.Now().Second()
	// Make the API request to upload metadata and file data.
	r, err := svc.Files.Insert(f).Media(m).Do()
	if err != nil {
		fmt.Printf("An error occurred uploading the document: %v\n", err)
	}
	tEnd := time.Now().Second()
	gd.totalSize += 1024 * 1024
	gd.totalSize += uint32(tEnd - tStart)
	fmt.Printf("Created: ID=%v, Title=%v\n", r.Id, r.Title)

	return r.Id
}

func (gd *GoogleDrive) CheckSize() (error, bool) {
    var c oauth.CacheFile = oauth.CacheFile(gd.authString)
    tok, err := c.Token()
    fmt.Println(err)
	t := &oauth.Transport{
		Config:    config,
		Transport: http.DefaultTransport,
        Token:tok,
	}
	// Create a new authorized Drive client.
	d, err := drive.New(t.Client())
	if err != nil {
		fmt.Printf("An error occurred creating Drive client: %v\n", err)
	}
	a, err := d.About.Get().Do()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
	}
	return err, a.QuotaBytesTotal-a.QuotaBytesUsed >= 1024*1024
}
