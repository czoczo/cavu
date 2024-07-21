// implements HTTP requests (heathcheck + )

package main

import (
	"io"
	"os"
	"net/url"
	"strings"
	"net/http"
	"path/filepath"
	log "github.com/sirupsen/logrus"
)


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