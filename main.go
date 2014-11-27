package main

import (
	"flag"
	"github.com/kvisscher/hollow-moose/slack/votes"
	"log"
	"net/http"
)

func main() {
	var slackAuthToken string
	var targetChannel string

	flag.StringVar(&slackAuthToken, "t", "", "The slack authentication token")
	flag.StringVar(&targetChannel, "c", "", "The channel to send and receive messages for")
	flag.Parse()

	if slackAuthToken == "" {
		log.Fatal("Invalid authentication token")
	}

	http.Handle("/slack", votes.New(slackAuthToken, targetChannel))

	server := &http.Server{Addr: ":8080"}

	log.Fatal(server.ListenAndServe())
}
