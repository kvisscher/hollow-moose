package votes

import (
	"encoding/json"
	"errors"
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
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	channel := request.FormValue("channel_name")

	if channel != "" && channel != handler.Channel {
		log.Println("Skipping message from channel", channel, "because it does not match the target", handler.Channel)
		return
	}

	trigger := request.FormValue("trigger_word")
	text := request.FormValue("text")

	response, err := handler.handleTrigger(trigger, text, request)

	if err != nil {
		log.Println("Error while handling trigger", trigger, err)
		return
	}

	bytes, err := json.Marshal(response)

	if err != nil {
		log.Println("Error encoding response")
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Write(bytes)
}

func (handler *VotesSlackHandler) handleTrigger(trigger string, text string, request *http.Request) (slack.Response, error) {
	var response slack.Response
	var err error

	if strings.Index(trigger, CommandPlusOne) == 0 {
		log.Println("Plus one triggered for", handler.CurrentVoteTarget)

		// Add 1 to the total score
		response = handler.handleUpdateScore(1)
	} else if strings.Index(trigger, CommandMinusOne) == 0 {
		log.Println("Minus one triggered for", handler.CurrentVoteTarget)

		// Remove 1 from the total score
		response = handler.handleUpdateScore(-1)
	} else if strings.Index(trigger, CommandVote) == 0 {
		// Retrieve the target that is going to be voted for from the text
		handler.CurrentVoteTarget = strings.TrimSpace(strings.Replace(text, trigger, "", 1))

		if _, ok := handler.Votes[handler.CurrentVoteTarget]; !ok {
			// Create new vote entry
			handler.Votes[handler.CurrentVoteTarget] = &UserScore{User: request.FormValue("user_name"), Votes: 0}

			log.Println("Added new entry for", handler.CurrentVoteTarget)
		}
	} else if strings.Index(trigger, CommandStats) == 0 {
		// TODO Post a message to show the stats of a user
	} else {
		err = errors.New(fmt.Sprintf("Unknown command %s", trigger))
	}

	return response, err
}

func (handler *VotesSlackHandler) handleUpdateScore(change int) slack.Response {
	var response slack.Response

	if score, ok := handler.Votes[handler.CurrentVoteTarget]; ok {
		score.Votes += change

		response = createResponse(fmt.Sprintf("Added +1 to score for %s new score is %d", handler.CurrentVoteTarget, score.Votes))
	}

	return response
}

func createResponse(text string) slack.Response {
	return slack.Response{Text: text}
}
