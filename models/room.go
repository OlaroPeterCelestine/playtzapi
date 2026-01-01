package models

// Room represents a music room/genre
type Room struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Genre       string `json:"genre"`
	Description string `json:"description"`
	Gradient    string `json:"gradient"`
	TextColor   string `json:"text_color"`
	Image       string `json:"image,omitempty"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}
