package main

import (
    "github.com/elastic/go-elasticsearch/v7"
    "log"
)

func main1() {
    es, err := elasticsearch.NewDefaultClient()
    if err != nil {
   	 log.Fatalf("Error creating the client: %s", err)
    }
    log.Println(elasticsearch.Version)

    res, err := es.Info()
    if err != nil {
   	 log.Fatalf("Error getting response: %s", err)
    }
    defer res.Body.Close()
    log.Println(res)
}