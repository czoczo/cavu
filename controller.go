// TO DELETE?

package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	generatedAvatarsPath  = "./avatars"
	downloadedAvatarsPath = "./downloadedAvatars"
	staticApiPath         = "/api/v1"
)

func firstWord(input string) string {
	// Split the input string by spaces
	input = strings.ToLower(input)
	words := strings.Fields(input)

	// Check if there are no spaces
	if len(words) <= 1 {
		// If no spaces, return the original string
		return input
	}

	// If there are spaces, return only the first word
	return words[0]
}

func downloadIcon(fullURLFile string) string {

	// create avatars folder if not existant
	log.Debug("Creating folder if non-existent: ", compiledVuePath+"/"+generatedAvatarsPath)
	os.MkdirAll(compiledVuePath+"/"+generatedAvatarsPath, os.ModePerm)

	// return if avatar generated locally
	if strings.HasPrefix(fullURLFile, generatedAvatarsPath) {
		return fullURLFile
	}

	// save to file if SVF data in string
	if strings.HasPrefix(fullURLFile, "data:image/svg+xml,") {
		// build avatar filename and path
		fileName := downloadedAvatarsPath + "/" + strToSha256(fullURLFile) + ".svg"
		writeStringFile(compiledVuePath+"/"+fileName, fullURLFile[20:])
		return fileName
	}

	// Validate URL
	_, err := url.Parse(fullURLFile)
	if err != nil {
		log.Fatal("Error validating URL: ", err)
	}

	// Build filename
	extension := filepath.Ext(fullURLFile)
	extension = strings.Split(extension, "?")[0]

	// create downloadedAvatars folder if not existant
	log.Debug("Creating folder if non-existent: ", compiledVuePath+"/"+downloadedAvatarsPath)
	os.MkdirAll(compiledVuePath+"/"+downloadedAvatarsPath, os.ModePerm)

	// build avatar filename and path
	fileName := downloadedAvatarsPath + "/" + strToSha256(fullURLFile) + extension

	// Create a blank file
	file, err := os.Create(compiledVuePath + "/" + fileName)
	if err != nil {
		log.Fatal("Error creating file: ", err)
	}
	defer file.Close()

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	// Download content and save to file
	resp, err := client.Get(fullURLFile)
	if err != nil {
		log.Fatal("Error downloading content: ", err)
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)
	log.Debug("Downloaded file %s with size %d\n", compiledVuePath+"/"+fileName, size)
	return fileName
}

func refreshItems() {
	for _, name := range dashboardItems.getKeys() {
		if *staticMode {
			wg.Add(1)
		}
		go crawlItem(name)
	}
	if *staticMode {
		wg.Wait()

		// create folder for static api file if not existant
		log.Debug("Creating folder if non-existent: ", filepath.Dir(compiledVuePath+"/"+staticApiPath))
		os.MkdirAll(filepath.Dir(compiledVuePath+staticApiPath), os.ModePerm)

		err := exportConfigAsJSONFile(dashboardItems.get(), compiledVuePath+staticApiPath)
		if err != nil {
			fmt.Printf("Error creating JSON file: %v\n", err)
		}
	}
}

func entriesApiHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Api request: ", r.Method, " ", r.URL.Path)
	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Marshal the dashboardItems to JSON
	jsonData, err := json.Marshal(dashboardItems.get())
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	// Write the JSON data to the response writer
	w.Write(jsonData)
}
