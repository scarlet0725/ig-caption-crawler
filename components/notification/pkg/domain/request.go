package domain

type NotificationRequest struct {
	IconURL string `json:"icon_url,omitempty"`
	Message string `json:"message,omitempty"`
	Name    string `json:"name,omitempty"`
}

type NotificationResponse struct {
	Ok bool `json:"ok,omitempty"`
}

type DiscordNotificationRequest struct {
	ProfileURL string `json:"profile_url,omitempty"`
	PostURL    string `json:"post_url,omitempty"`
	Title      string `json:"title,omitempty"`
	NotificationRequest
}
