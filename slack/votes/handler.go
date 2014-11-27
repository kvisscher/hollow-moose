package votes

import (
	"encoding/json"
	"fmt"
	"github.com/kvisscher/hollow-moose/slack"
	"log"
	"net/http"
	"strings"
)

const (
	CommandPlusOne  = "+1"
	CommandMinusOne = "-1"
	CommandVote     = "vote"
	CommandStats    = "stats"
)

type VotesSlackHandler struct {
	Token             string
	Channel           string
	Votes             map[string]*UserScore
	CurrentVoteTarget string
}

type UserScore struct {
	User  string // User that posted the url
	Votes int    // The number of votes the url received
}

func New(token string, channel string) *VotesSlackHandler {
	return &VotesSlackHandler{
		Token:   token,
		Channel: channel,
		Votes:   make(map[string]*UserScore),
	}
}

// Handler for the external webhook of Slack
func (handler *VotesSlackHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	token := request.FormValue("token")

	if token != handler.Token {
		log.Println("Received invalid token. Ignoring this request")
		return
	}

	channel := request.FormValue("channel_name")

	if channel != handler.Channel {
		log.Println("Skipping message from channel", channel, "because it does not match the target", handler.Channel)
		return
	}

	trigger := request.FormValue("trigger_word")
	text := request.FormValue("text")

	if strings.Index(trigger, CommandPlusOne) == 0 {
		log.Println("Plus one triggered for", handler.CurrentVoteTarget)
		if score, ok := handler.Votes[handler.CurrentVoteTarget]; ok {
			score.Votes += 1

			log.Println("Added +1 to score for", handler.CurrentVoteTarget, "new score is", score.Votes)

			respondToSlack(fmt.Sprintf("Added +1 to score for %s new score is %d", handler.CurrentVoteTarget, score.Votes), writer)
		}
	} else if strings.Index(trigger, CommandMinusOne) == 0 {
		// Subtract one from the score
		if score, ok := handler.Votes[handler.CurrentVoteTarget]; ok {
			score.Votes -= 1

			log.Println("Removed -1 from score for", handler.CurrentVoteTarget, "new score is", score.Votes)

			respondToSlack(fmt.Sprintf("Removed -1 from score for %s new score is %d", handler.CurrentVoteTarget, score.Votes), writer)
		}
	} else if strings.Index(trigger, CommandVote) == 0 {
		// Retrieve the url from the text
		handler.CurrentVoteTarget = strings.TrimSpace(strings.Replace(text, trigger, "", 1))

		if _, ok := handler.Votes[handler.CurrentVoteTarget]; !ok {
			// Create new entry
			handler.Votes[handler.CurrentVoteTarget] = &UserScore{User: request.FormValue("user_name"), Votes: 0}

			log.Println("Added new entry for", handler.CurrentVoteTarget)
		}
	} else if strings.Index(trigger, CommandStats) == 0 {
		// TODO Post a message to show the stats of a user
	} else {
		log.Println("Unknown command", trigger)
	}
}

func respondToSlack(text string, writer http.ResponseWriter) {
  bytes, err := json.Marshal(slack.Response{Text: text})

  if err != nil {
    log.Println("Error encoding response")
    return
  }

  writer.Header().Set("Content-Type", "application/json; charset=utf-8")
  writer.Write(bytes)
}
