package main

import (
	"encoding/json"
	"flag"
	"fmt"
	heimdall "github.com/gojektech/heimdall/httpclient"
	"log"
	"net/http"
	"time"
)

const discordURL string = "https://discordapp.com/api/v6"

type message struct {
	ID      string `json:"id"`
	Author  author `json:"author"`
	Content string `json:"content"`
}

type author struct {
	Username string `json:"username"`
}

func main() {
	channel := flag.String("channel", "", "Discord channel or DM to erase.")
	authorizationHeader := flag.String(
		"authorizationHeader",
		"",
		"Discord Authorization header to use when issuing Discord API requests.",
	)
	username := flag.String(
		"username",
		"",
		"(optional) Discord username of the logged in user",
	)

	flag.Parse()

	if *channel == "" {
		log.Fatalf("Provide a channel using -channel to continue.")
	}

	if *authorizationHeader == "" {
		log.Fatalf(
			"Please provide the authorization header to use using " +
				"-authorizationHeader to continue.",
		)
	}

	log.Printf("Starting deletion process for channel/message %s", *channel)

	client := heimdall.NewClient()

	headers := http.Header{}
	headers.Set("authorization", *authorizationHeader)

	var before string

	for {
		fetchURL := "%s/channels/%s/messages?limit=25"

		if before != "" {
			fetchURL += "&before=" + before
			log.Printf("Fetching messages before %s", before)
		} else {
			log.Printf("Fetching messages")
		}

		res, err := client.Get(
			fmt.Sprintf(fetchURL, discordURL, *channel),
			headers,
		)

		if err != nil {
			log.Fatalf("Error thrown when fetching messages: %+v", err)
		}

		var messages []message
		err = json.NewDecoder(res.Body).Decode(&messages)

		if len(messages) == 0 {
			break
		}

		log.Printf("%d messages fetched, preparing for deletion", len(messages))

		for _, message := range messages {
			if *username != "" && message.Author.Username != *username {
				continue
			}

			deletionURL := "%s/channels/%s/messages/%s"
			log.Printf("Deleting message %s", message.ID)

			res, err := client.Delete(
				fmt.Sprintf(deletionURL, discordURL, *channel, message.ID),
				headers,
			)

			if err != nil {
				log.Printf("Unable to delete message %s: %+v", message.ID, err)
				continue
			} else if res.StatusCode >= 400 {
				log.Printf(
					"Error code when deleting message %s: %d",
					message.ID,
					res.StatusCode,
				)
				continue
			}

			log.Printf("Message %s deleted", message.ID)
			time.Sleep(250 * time.Millisecond)
		}

		before = messages[len(messages)-1].ID
		time.Sleep(250 * time.Millisecond)
	}

	log.Printf("Finished deleting all messages")
}
