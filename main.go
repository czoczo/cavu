// application starting point

package main

import (
	"flag"
	"net/http"
	"time"
	"sync"
)

const (
	version         = "0.0.65"
	configFilePath  = "./config/main.yaml"
	itemsFilePath   = "./config/items.yaml"
	staticFilesPath = "./frontend"
	compiledVuePath = staticFilesPath + "/dist"
	sourceVuePath   = staticFilesPath + "/src"
)

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
