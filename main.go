package main

import (
	"log"
	"net/http"
	"os"

	"news-graphql/graphql"
	"news-graphql/news"
)

func main() {
	apiKey := os.Getenv("NEWSAPI_KEY")
	if apiKey == "" {
		log.Fatal("NEWSAPI_KEY environment variable required")
	}

	client := news.NewNewsClient(apiKey)

	schema := graphql.NewSchema(client)

	http.Handle("/query", schema.Handler())
	// Serve index.html for quick demo
	http.Handle("/", http.FileServer(http.Dir("./")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
