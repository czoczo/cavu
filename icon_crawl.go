// implement order of icon sources to crawl

package main

import (
	log "github.com/sirupsen/logrus"
	"strings"
)

func crawlItem(name string) {

	// in static generation mode wait,
	// for all items to complete before exiting
	if *staticMode {
		defer wg.Done()
	}

	dashboardItem, ok := dashboardItems.read(name)
	if !ok {
		return
	}

	// uncomment for debuging single item
	//if name != "PoeticMetric" { return }

	// echo metadata
	log.Info("Starting icon crawl for '", name, "' at '", dashboardItem.URL)

	findHtmlTitle(name, &dashboardItem)

	for key, val := range dashboardItem.Labels {
		// check for app.kubernetes.io/instance label
		if key == "app.kubernetes.io/instance" {
			findIconGitHub(&dashboardItem, val)
			log.Debug("GitHub icon search, based on '", val, "' returned: ", dashboardItem.IconURL)
		}
		// check for app.kubernetes.io/name label
		if key == "app.kubernetes.io/name" {
			findIconGitHub(&dashboardItem, val)
			log.Debug("GitHub icon search, based on '", val, "' returned: ", dashboardItem.IconURL)
		}
	}
	checkedNames := map[string]bool{}
	
	// check Icons on GitHub based on Ingress name
	checkedNames[name] = true
	findIconGitHub(&dashboardItem, strings.ToLower(name))

	// check Icons on GitHub based on site title first word
	titleFirstWord := firstWord(dashboardItem.WebpageTitle)
	if ! checkedNames[titleFirstWord] {
		checkedNames[titleFirstWord] = true
		findIconGitHub(&dashboardItem, titleFirstWord)
	}

	// check Icons on GitHub based on site title spaces to dashes
	titleWithDashes := strings.ToLower(strings.ReplaceAll(dashboardItem.WebpageTitle, " ", "-"))
	if ! checkedNames[titleWithDashes] {
		checkedNames[titleWithDashes] = true
		findIconGitHub(&dashboardItem, titleWithDashes)
	}

	// check header for PNG
	findHtmlIcon(&dashboardItem, "png")

	// check header for SVG
	findHtmlIcon(&dashboardItem, "svg")

	// check header with https://pkg.go.dev/go.deanishe.net/favicon
	findHtmlIconDeanishe(&dashboardItem)

	// check for first level of DNS domain
	addressPrefix := strings.Split(getHostFromURL(dashboardItem.URL), ".")[0]
	if ! checkedNames[addressPrefix] {
		checkedNames[addressPrefix] = true
		findIconGitHub(&dashboardItem, addressPrefix)
	}

	// last resort - generate avatar
	getGeneratedIcon(&dashboardItem, name)

	// download if static mode
	if *staticMode {
		downloadedIconFile := downloadIcon(dashboardItem.IconURL)
		log.Info("Downloaded icon file: ", downloadedIconFile, ", from: ", dashboardItem.IconURL)
		dashboardItem.IconURL = downloadedIconFile
	}

	// write result
	dashboardItems.write(name, dashboardItem)
	log.Info("Icon crawl result for '", name, "': icon - ", dashboardItem.IconURL, ", title - ", dashboardItem.WebpageTitle)
}
