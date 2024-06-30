// application starting point

package main

import (
	_ "embed"
	"flag"
	"net/http"
	"time"
	"sync"
)

const (
	configFilePath  = "./config/main.yaml"
	itemsFilePath   = "./config/items.yaml"
	staticFilesPath = "./frontend"
	compiledVuePath = staticFilesPath + "/dist"
	sourceVuePath   = staticFilesPath + "/src"
)

//go:embed VERSION_APP.txt
var version string

var config Config
var staticItems StaticItems
var staticMode *bool
var wg sync.WaitGroup
var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}


func main() {
	staticMode = flag.Bool("static", false, "Single shot static content dashboard generation.")
	flag.Parse()

	dashboardItems.items = make(map[string]DashEntry)

	// config.go
	loadConfig()

	// customisation.go
	updateFrontendFiles()

	if *staticMode {
		return
	}

	// kubernetes.go
	go getAndWatchKubernetesIngressItems()

	// httpserver.go
	initHttpServer()
}
