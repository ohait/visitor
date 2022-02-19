package data

import "time"

type Person struct {
	Id      string  `json:"id"`
	Email   string  `json:"email,omitempty"`
	Details Details `json:"details"`
}

type Details struct {
	Name      string   `json:"name,omitempty"`
	Addresses []string `json:"address"`
	Phone     []string `json:"phone"`
}

type Event struct {
	Time     time.Time `json:"time"`
	Person   string    `json:"person"` // Ref Person.Id
	Type     EventType `json:"type"`   // "landing", "video"
	Campaign string    `json:"campaign,omitempty"`
	Url      string    `json:"url,omitempty"`
}

type EventType string

const (
	EventTypeLandingPage EventType = "landing"
	EventTypeVideoPlay   EventType = "video"
)
