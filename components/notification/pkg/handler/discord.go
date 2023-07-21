package handler

import (
	"bytes"
	"encoding/json"
	"github.com/scarlet0725/insta-caption-crawler/components/notification/pkg/domain"
	"net/http"
)

type DiscordNotificationHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type discordNotificationHandler struct {
	webhookEndpoint string
}

func NewDiscordNotificationHandler(webhookEndpoint string) DiscordNotificationHandler {
	return &discordNotificationHandler{
		webhookEndpoint: webhookEndpoint,
	}
}

func (d *discordNotificationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := domain.DiscordNotificationRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	embed := domain.Embed{
		Title:       request.Title,
		Description: request.Message,
		URL:         request.PostURL,
		Author: &domain.EmbedAuthor{
			Name:    request.Name,
			URL:     request.ProfileURL,
			IconURL: request.IconURL,
		},
	}

	webhookBody := domain.WebhookMessageCreate{
		AvatarURL: request.IconURL,
		Username:  request.Name,
		Embeds: []domain.Embed{
			embed,
		},
	}

	body, err := json.Marshal(webhookBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = http.Post(d.webhookEndpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}
