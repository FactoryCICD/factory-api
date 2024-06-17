package webhooks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type WebHookHandler interface {
	IncomingWebHook(w http.ResponseWriter, r *http.Request)
}

func Routes(handler WebHookHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", handler.IncomingWebHook)

	return router
}

type githubWebHookHandler struct{}

func NewGithubWebHookHandler() WebHookHandler {
	return &githubWebHookHandler{}
}

func (gwh *githubWebHookHandler) IncomingWebHook(w http.ResponseWriter, r *http.Request) {
	// Get the hook type from the header X-Github-Event
	event := r.Header.Get("X-GitHub-Event")
	switch event {
	case "ping":
		fmt.Println("Recieved a ping!")
	case "push":
		var pushEvent GitHubPushEvent
		err := json.NewDecoder(r.Body).Decode(&pushEvent)
		if err != nil {
			fmt.Println("Error decoding request body")
		}
		fmt.Println(pushEvent)
	default:
		fmt.Printf("%s not supported\n", event)
	}
	w.WriteHeader(http.StatusOK)
}

type GitHubPing struct {
	Zen string `json:"zen"`
}

// GitHubPushEvent describes a push event from github
// docs: https://docs.github.com/en/webhooks/webhook-events-and-payloads#push
type GitHubPushEvent struct {
	After      string        `json:"after"`
	BaseRef    string        `json:"base_ref"`
	Before     string        `json:"before"`
	Commits    []interface{} `json:"commits"`
	Pusher     GitHubPusher  `json:"pusher"`
	Ref        string        `json:"ref"`
	Repository interface{}   `json:"repository"`
	Sender     interface{}   `json:"sender"`
}

type GitHubPusher struct {
	Date     string `json:"date"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
}
