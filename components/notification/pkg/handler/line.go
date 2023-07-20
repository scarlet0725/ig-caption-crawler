package handler

import (
	"encoding/json"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/scarlet0725/insta-caption-crawler/components/notification/pkg/domain"
	"net/http"
)

type LineWebhookHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type LineNotificationHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type lineWebhookHandler struct {
	line *linebot.Client
}

type lineNotificationHandler struct {
	line *linebot.Client
}

func NewLineWebhookHandler(line *linebot.Client) LineWebhookHandler {
	return &lineWebhookHandler{
		line: line,
	}
}

func NewLineNotificationHandler(line *linebot.Client) LineNotificationHandler {
	return &lineNotificationHandler{
		line: line,
	}
}

func (l *lineWebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	events, err := l.line.ParseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(events) < 1 {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (l *lineNotificationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	request := domain.NotificationRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := linebot.NewTextMessage(request.Message)
	sender := linebot.NewSender(request.Name, request.IconURL)
	msg.WithSender(sender)

	ctx := r.Context()
	_, err := l.line.BroadcastMessage(msg).WithContext(ctx).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := domain.NotificationResponse{
		Ok: true,
	}

	_ = json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusOK)

}
