package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	initFlag := flag.Bool("init", false, "Initialize the config file with default values")
	deleteCache := flag.Bool("delete-cache", false, "Delete the cache file and exit")
	flag.Parse()

	if *deleteCache {
		err := os.Remove("cache.db")
		if err != nil {
			log.Fatalf("Error deleting cache file: %v", err)
		}
		fmt.Println("All cache has been deleted")
		return
	}

	if *initFlag {
		initializeConfig()
		return
	}

	config := loadConfig()

	// initDB()
	initCache(config.MyForwardProxy.ExpirationPeriod)

	mux := http.NewServeMux()

	mux.HandleFunc("/", proxyHandler)

	fmt.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
