package types

type Project struct {
	ID            string `json:"id"`
	URL           string `json:"url"`
	WebhookSecret string `json:"webhook_secret"`
}

type CreateProject struct {
	URL           string `json:"url"`
	WebhookSecret string `json:"webhook_secret"`
}
