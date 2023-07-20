package domain

type NotificationRequest struct {
	IconURL string `json:"icon_url,omitempty"`
	Message string `json:"message,omitempty"`
	Name    string `json:"name,omitempty"`
}

type NotificationResponse struct {
	Ok bool `json:"ok,omitempty"`
}
