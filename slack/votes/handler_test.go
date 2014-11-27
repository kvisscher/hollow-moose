package votes_test

import (
	"github.com/kvisscher/hollow-moose/slack/votes"
  // "github.com/kvisscher/hollow-moose/slack"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
  "io/ioutil"
  // "encoding/json"
  "log"
)

func TestThatInvalidTokenIsRejected(t *testing.T) {
	token := "my-token"
	handler := votes.New(token, "channel")

	server := httptest.NewServer(handler)
	defer server.Close()

  // Tokens  that are expected to fail to authenticate
  tokensToValidate := []string{"invalid-token", ""}

  for _, token := range tokensToValidate {
  	response, err := http.PostForm(server.URL, url.Values{
  		"token": {token},
  	})

    if err != nil {
      t.Fatal(err)
    }

    if response.StatusCode != http.StatusUnauthorized {
      t.Fatal("expected", http.StatusUnauthorized, "got", response.StatusCode)
    }
  }
}

func TestThatScoreIsSaved(t *testing.T) {
    token := "my-token"
    channel := "channel"
    userName := "user"
    voteTarget := "vote-target"

    handler := votes.New(token, channel)

    server := httptest.NewServer(handler)
    defer server.Close()

    response, err := http.PostForm(server.URL, url.Values{
            "token": {token},
            "channel_name": {channel},
            "user_name": {userName},
            "trigger_word": {votes.CommandVote},
            "text": {votes.CommandVote + " " + voteTarget},
      })

    _, err = ioutil.ReadAll(response.Body)

    response.Body.Close()

    if err != nil {
      t.Fatal("error reading response", err)
    }

    if response.StatusCode != http.StatusOK {
      t.Fatal("expected response to succeed, got status code", response.StatusCode)
    }

    if _, ok := handler.Votes[voteTarget]; !ok {
      t.Fatal("expected to have an entry for target", voteTarget)
    }
}

func init() {
  log.Println("init claled")
}

// func TestPlusCommandResponds(t *testing.T) {
//   token := "my-token"
//   channel := "channel"
//   userName := "user"
//
//   handler := votes.New(token, channel)
//
//   server := httptest.NewServer(handler)
//   defer server.Close()
//
//   response, err = http.PostForm(server.URL, url.Values{
//       "token": {token},
//       "channel_name": {channel},
//       "user_name": {userName},
//       "trigger_word"
//     })
// }
